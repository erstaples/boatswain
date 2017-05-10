// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/medbridge/mocking/factories"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var build = Build{}

type StagingValuesYAML struct {
	ImageTag string
}

// stageCmd represents the stage command
var stageCmd = &cobra.Command{
	Use:   "stage [push|delete]",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use boatswain stage [push|delete]")
	},
}

var stagePushCmd = &cobra.Command{
	Use:   "push [appnames] [domain]",
	Short: "Push an application(s) to staging",
	Long: `Push an application or bundle of applications to staging

	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStagePush(args)
	},
}

var stageDeleteCmd = &cobra.Command{
	Use:   "delete [appnames] [domain]",
	Short: "Delete an application(s) from staging",
	Long: `Delete an application or bundle of applications from staging
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStageDelete(args)
	},
}

func init() {
	RootCmd.AddCommand(stageCmd)
	stageCmd.AddCommand(stagePushCmd, stageDeleteCmd)

	targetDesc := "Target Dockerfiles to include in the build. Example: --targets medbridge-phpfpm builds Dockerfile.medbridge-phpfpm"
	stagePushCmd.Flags().StringVar(&build.Targets, "targets", "", targetDesc)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stageCmd.Flags().StringVar(&StageFlags."toggle", "t", false, "Help message for toggle")

}

func RunStagePush(args []string) {
	if len(args) != 2 {
		fmt.Printf("Unexpected number of args. Expected 2, got %c", len(args))
		return
	}
	appnames := args[0]
	domain := args[1]
	apps := strings.Split(appnames, ",")
	var selectedBuilds []Build
	config := Config{}
	configPath := viper.ConfigFileUsed()
	yamlFile, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(yamlFile, &config)

	for _, appname := range apps {
		//get build object from config corresponding to appname
		for _, build := range config.Builds {
			if appname == build.Name {
				selectedBuilds = append(selectedBuilds, build)
			}
		}
	}

	for _, build := range selectedBuilds {
		stagingYaml := StagingValuesYAML{}
		stagingYaml.ImageTag = runBuild(build, domain)
		yaml := getStagingYaml(stagingYaml)
		valuesPath := generateValuesFile(build, config, domain, yaml)
		runRelease(build, valuesPath, domain)
		genIngress(build, domain)
	}

	/**

	Needs to:
		* parse comma-delimited list of appnames
		* call the build.sh script
		* accept and pass in an optional list of build targets
		* take the resulting commit sha and build a new values file w/ pattern: values.staging.<domain>.yaml
		* call release command set to staging and pass in new staging values file
		* accept a --db option, update helm templates so that it does (or not) provisions db
		* deploy to test namespace (maybe accept that as an option?)
		* create necessary ingresses
		* correctly wire up services for appnames
	**/
}

func RunStageDelete(args []string) {

}

func runRelease(build Build, valuesFile string, domain string) {

	args := []string{build.Name}
	options := ReleaseOptions{
		Environment:       "staging",
		DryRun:            false,
		Namespace:         "default", //todo: change
		Packfile:          valuesFile,
		Xdebug:            false,
		NoExecute:         false, //todo: false
		PackageIDOverride: domain,
	}
	RunRelease(args, options)
}

func runBuild(build Build, domain string) string {
	// cmdFactory := factories.CommandFactory{}
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
	cmd := exec.Command(cmdName, cmdArgs...)

	//https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("build.sh | %s\n", scanner.Text())
		}
	}()

	//TODO: use factory

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	sha := getGitCommitSha(build)
	return string(sha[:])
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

func printWrapper(msg string) {
	fmt.Printf("boatswain | %s\n", msg)
}

func getStagingYaml(yaml StagingValuesYAML) string {
	tmpl, err := template.New("values").Parse(`

##################################################################
##                                                              ##
##  Autogenerated file. Changes made here will be overwritten.  ##
##                                                              ##
##################################################################

Boatswain:
  ImageTag: "{{.ImageTag}}"

`)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, yaml)
	s := doc.String()

	if err != nil {
		panic(err)
	}
	return s
}

func generateValuesFile(build Build, config Config, domain string, yaml string) string {
	fileName := "values.staging." + domain + ".yaml"
	valuesPath := config.Release + "/" + build.Name + "/autogenerated/" + fileName
	printWrapper("Creating : " + fileName)
	printWrapper("Values path: " + valuesPath)

	err := ioutil.WriteFile(valuesPath, []byte(yaml), 0777)
	if err != nil {
		panic(err)
	}
	return valuesPath
}

func genIngress(build Build, domain string) {
	cmdFactory := &factories.CommandFactory{}
	args := []string{domain + ".k8staging.medbridgeeducation.com"}
	options := GenIngressFlags{
		Service:     domain + "-medbridge", //todo: make less dumb... maybe config?? //build needs a ingress_service
		EnableTLS:   false,
		ServicePort: "80",
	}
	RunGenIngress(args, cmdFactory, options)
}
