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

	yaml "gopkg.in/yaml.v2"

	"github.com/medbridge/boatswain/lib"
	"github.com/medbridge/mocking/factories"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var build = lib.Build{}
var branchName string
var serviceMapName string
var serviceMapConfig ServiceMapConfig
var serviceMap ServiceMap
var configMapEntry lib.StagingConfigMapEntry
var stagingConfigMap lib.StagingConfigMap

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
	stagingConfigMap.LoadConfigMap()
	configMapEntry.Name = branchName
	configMapEntry.Ingress = branchName + ".k8staging.medbridgeeducation.com"

	config := Config{}
	configPath := viper.ConfigFileUsed()
	yamlFile, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(yamlFile, &config)

	selectedBuilds := getBuilds(config)

	cfTemplate := lib.CloudFormationTemplate{Output: make(map[string]string)}

	if len(serviceMap.CloudFormationTemplate) > 0 {
		cfTemplate = *lib.NewCloudFormationTemplate(serviceMap.CloudFormationTemplate)
		cfTemplate.CreateStack(branchName)
		configMapEntry.CloudFormationStack = cfTemplate.StackName
	}

	env := convertMapToEnvVars(serviceMap)
	imageTags := make(map[string]string)

	for _, build := range selectedBuilds {

		fmt.Printf("\nRunning build %s", build.Name)
		imageTags[build.Name] = build.Exec()

	}

	for _, svc := range serviceMap.Test {
		values := lib.NewValues(branchName, svc, imageTags[svc], env)
		values.CloudFormationValues = cfTemplate.Output
		runRelease(svc, values.Write())
	}
	genIngress()

	stagingConfigMap.AddConfig(configMapEntry)
}

func runRelease(name string, valuesFile string) {

	args := []string{name}
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
	configMapEntry.HelmDeployments = append(configMapEntry.HelmDeployments, name)
}

func genIngress() {
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

func getBuilds(config Config) []lib.Build {
	var builds []lib.Build
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
