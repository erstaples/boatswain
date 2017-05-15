// Copyright © 2017 MedBridge Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/medbridge/boatswain/utilities"
	"github.com/medbridge/mocking/factories"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var build = Build{}
var branchName string
var serviceMapName string
var serviceMapConfig ServiceMapConfig
var serviceMap ServiceMap

var stagePushCmd = &cobra.Command{
	Use:   "push [appnames] [domain]",
	Short: "Push an application(s) to staging",
	Long: `Push an application or bundle of applications to staging

	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStagePush(args)
	},
}

func init() {
	StageCmd.AddCommand(stagePushCmd)
}

func RunStagePush(args []string) {
	if len(args) != 2 {
		fmt.Printf("Unexpected number of args. Expected 2, got %c", len(args))
		return
	}

	serviceMapName = args[0]
	branchName = args[1]
	loadServiceMap()

	config := Config{}
	configPath := viper.ConfigFileUsed()
	yamlFile, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(yamlFile, &config)

	selectedBuilds := getSelectedBuilds(config)

	for _, build := range selectedBuilds {
		fmt.Printf("\n running build %s", build)
		imageTag := runBuild(build)

		env := convertMapToEnvVars(serviceMap)

		stagingYaml := StagingValuesYAML{
			ImageTag: imageTag,
			Env:      env,
		}

		if len(serviceMap.CloudFormationTemplate) > 0 {
			cfValues := readCloudFormationTemplate()
			stagingYaml.CloudFormationValues = cfValues
		}

		yaml := getStagingYaml(stagingYaml)

		fmt.Printf("%s", yaml)
		valuesPath := createValuesFile(build, config, yaml)
		runRelease(build, valuesPath)
		genIngress(build)
	}
}

func runBuild(build Build) string {
	cmdName := "/bin/bash"
	cmdArgs := []string{build.Path, "push"}

	targetsString := build.Targets
	targets := strings.Split(targetsString, ",")

	if len(targets) > 0 {
		for _, target := range targets {
			cmdArgs = append(cmdArgs, target)
		}
	}
	os.Chdir(build.Rootpath)
	utilities.ExecStreamOut(cmdName, cmdArgs, "build.sh")

	sha := getGitCommitSha(build)
	return string(sha[:])
}

func runRelease(build Build, valuesFile string) {

	args := []string{build.Name}
	options := ReleaseOptions{
		Environment:       "staging",
		DryRun:            false,
		Namespace:         "default",
		Packfile:          valuesFile,
		Xdebug:            false,
		NoExecute:         false,
		PackageIDOverride: branchName,
	}
	RunRelease(args, options)
}

func getGitCommitSha(build Build) []byte {
	os.Chdir(build.Rootpath)
	cmdName := "git"
	cmdArgs := []string{"show", "-s", "--pretty=format:%h"}
	out, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		panic(err)
	}
	return out
}

func getStagingYaml(yaml StagingValuesYAML) string {
	funcs := template.FuncMap{
		"ToEnv": func(s string) string {
			s = strings.Replace(s, "-", "_", -1)
			return strings.ToUpper(s) + "_HOST"
		},
	}
	tmpl, err := template.New("values").Funcs(funcs).Parse(`

##################################################################
##                                                              ##
##  Autogenerated file. Changes made here will be overwritten.  ##
##                                                              ##
##################################################################

Boatswain:
  ImageTag: "{{.ImageTag}}"
  ServiceEnv:
    {{- range $key, $value := .Env }}
    {{ $key | ToEnv }}: {{ $value }}
	{{- end }}
  CloudFormationValues:
    {{- range $key, $value := .CloudFormationValues }}
    {{ $key }}: {{ $value }}
    {{- end }}

`)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, yaml)
	s := doc.String()

	if err != nil {
		panic(err)
	}
	return s
}

func createValuesFile(build Build, config Config, yaml string) string {
	fileMode := os.FileMode(0777)
	fileName := "values.staging." + branchName + ".yaml"
	autogenPath := config.Release + "/" + build.Name + "/autogenerated/"
	if !pathExists(autogenPath) {
		os.Mkdir(autogenPath, fileMode)
	}
	valuesPath := autogenPath + fileName
	utilities.PrintWrapper("boatswain", "Creating : "+fileName)
	utilities.PrintWrapper("boatswain", "Created values file: "+valuesPath)

	err := ioutil.WriteFile(valuesPath, []byte(yaml), fileMode)
	if err != nil {
		panic(err)
	}
	return valuesPath
}

func genIngress(build Build) {
	cmdFactory := &factories.CommandFactory{}

	tmpl, _ := template.New("ingress_host").Parse(serviceMapConfig.Ingress.Template)
	var doc bytes.Buffer
	ingressName := struct {
		BranchName string
	}{branchName}
	err := tmpl.Execute(&doc, ingressName)
	if err != nil {
		panic(err)
	}
	ingressHost := doc.String()
	args := []string{ingressHost}

	options := GenIngressFlags{
		Service:     branchName + "-" + serviceMapConfig.Ingress.Service,
		EnableTLS:   false,
		ServicePort: serviceMapConfig.Ingress.Port,
	}

	RunGenIngress(args, cmdFactory, options)
}

func loadServiceMap() {
	//get the service serviceMap file
	path := viper.GetString("release") //TODO: capitalize this
	fullPath := path + "/.servicemap/staging.yaml"
	valuesBytes, err := ioutil.ReadFile(fullPath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(valuesBytes, &serviceMapConfig)
	if err != nil {
		panic(err)
	}

	for _, sMap := range serviceMapConfig.ServiceMaps {
		if serviceMapName == sMap.Name {
			serviceMap = sMap
		}
	}
}

func convertMapToEnvVars(serviceMap ServiceMap) map[string]string {
	env := make(map[string]string)

	for _, testSvc := range serviceMap.Test {
		env[testSvc] = branchName + "-" + testSvc
	}

	for _, stagingSvc := range serviceMap.Staging {
		env[stagingSvc] = "staging-" + stagingSvc
	}

	return env
}

func getSelectedBuilds(config Config) []Build {
	var builds []Build
	fmt.Printf("%s", serviceMap)
	for _, testSvc := range serviceMap.Test {
		for _, build := range config.Builds {
			if testSvc == build.Name {
				builds = append(builds, build)
			}
		}
	}
	return builds
}

func readCloudFormationTemplate() map[string]string {
	path := viper.GetString("release")
	cf := serviceMap.CloudFormationTemplate
	cfTemplate := path + "/.cloudformation/" + cf + ".yaml"
	yamlBytes, err := ioutil.ReadFile(cfTemplate)
	if err != nil {
		panic(err)
	}
	return loadCloudFormationTemplate(cf, yamlBytes)
}

func loadCloudFormationTemplate(cf string, yamlBytes []byte) map[string]string {
	fmt.Printf("\nRunning CloudFormation stack [%s]", cf)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	svc := cloudformation.New(sess)

	templateBody := string(yamlBytes)
	stackName := cf + "-" + branchName
	cloudFormationValues := make(map[string]string)

	params := &cloudformation.CreateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
	}

	describeStacksParams := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	out, err := svc.CreateStack(params)
	if err != nil {
		if strings.Contains(err.Error(), "AlreadyExistsException") {
			fmt.Printf("\nCloudFormation stack [%s] already exists. Skipping...", cf)
			descOut, err := svc.DescribeStacks(describeStacksParams)
			if err != nil {
				panic(err)
			}
			cloudFormationValues = parseOutput(descOut, cloudFormationValues)
			return cloudFormationValues
		}
		fmt.Printf("%s", err)
		panic(err)
	} else {
		fmt.Printf("%s", out)
	}

	stackReady := false

	for stackReady != true {

		descOut, err := svc.DescribeStacks(describeStacksParams)
		if err != nil {
			fmt.Printf("%s", err)
			panic(err)
		} else {
			fmt.Printf("\nCloudFormation stack [%s] is creating...", cf)
		}

		if *descOut.Stacks[0].StackStatus == "CREATE_COMPLETE" {
			stackReady = true
			fmt.Printf("\nCloudFormation stack [%s] ready...\n", cf)
			cloudFormationValues = parseOutput(descOut, cloudFormationValues)
		}

		time.Sleep(time.Second * 7)
	}

	return cloudFormationValues
}

func parseOutput(descOut *cloudformation.DescribeStacksOutput, cloudFormationValues map[string]string) map[string]string {
	stack := descOut.Stacks[0]
	for _, cfOutput := range stack.Outputs {
		trimKey := strings.TrimSpace(*cfOutput.OutputKey)
		trimVal := strings.TrimSpace(*cfOutput.OutputValue)
		cloudFormationValues[trimKey] = trimVal
	}
	return cloudFormationValues
}
