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
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var colorNone = "\033[00m"
var colorYellow = "\033[01;33m"
var colorGreen = "\033[01;32m"

var env string
var dryrun bool
var ns string
var packfile string
var xdebug bool
var noExecute bool
var packageIDOverride string
var optSetValues string

type Config struct {
	ReleasePath string
}

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <appName>",
	Short: "Deploys application to minikube, staging, or production cluster",
	Long: `Wraps around a helm install command to automate common helm configuration
	options. Sets packageId, environment, and other important values.Execute
	
	Examples:
	boatswain release medbridge -x
	Release medbridge in the minikube cluster with XDebug enabled

	release medbridge -e staging
	Release medbridge in the staging cluster

	release medbridge -e staging -n test
	Release medbridge in the staging cluster, test namespace
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Required argument: releaseName")
			return
		}

		releaseName := args[0]

		var xdebugHost string
		var fullReleaseName string
		var packageId string

		environments := []string{"development", "dev", "staging", "stage", "production", "prod"}

		// based on environment, set default packageId, packfile, and context
		switch env {
		case environments[0], environments[1]:
			if len(packageId) == 0 {
				packageId = "dev"
			}
			if len(packfile) == 0 {
				packfile = "values.env.yaml"
			}
			useK8sCurrContext("minikube")

		case environments[2], environments[3]:
			if len(packageId) == 0 {
				packageId = "staging"
			}
			if len(packfile) == 0 {
				packfile = "values.staging.yaml"
			}
			useK8sCurrContext("staging")
		case environments[4], environments[5]:
			if len(packageId) == 0 {
				packageId = "prod"
			}
			if len(packfile) == 0 {
				packfile = "values.prod.yaml"
			}
			useK8sCurrContext("production")
		default:
			fmt.Println("Invalid environment: " + env)
			os.Exit(1)
		}

		if len(packageIDOverride) > 0 {
			packageId = packageIDOverride
		}

		fullReleaseName = packageId + "-" + releaseName

		//xdebug option turned on, so get the xdebug host ip address
		if xdebug {

			var (
				cmdOut []byte
				err    error
			)
			if cmdOut, err = getXDebugHost(); err != nil {
				fmt.Fprintln(os.Stderr, "There was an error running ipconfig command: ", err)
				os.Exit(1)
			}
			xdebugHost = string(cmdOut)
			fmt.Println("Output: ", xdebugHost)
		}

		releasePath := viper.GetString("ReleasePath")
		appPath := releasePath + "/" + releaseName

		fmt.Printf("Deploying: %s\n", appPath)

		//build helm cmd
		//add standard values to be made available in the helm releases
		timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		setValues := "environment=" + env + ",packageId=" + packageId + ",timestamp='" + timestamp + "'"
		if xdebug {
			setValues += ",xdebugHost=" + xdebugHost
		}

		//fully qualified path
		packfileFullPath := appPath + "/" + packfile

		execHelmUpgradeCmd(fullReleaseName, appPath, setValues, packfileFullPath, packfile, ns)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)

	//set option flags
	releaseCmd.Flags().StringVarP(&env, "environment", "e", "development", "Target environment for the release. 'production', 'staging', and 'development' are valid options")
	releaseCmd.Flags().BoolVarP(&dryrun, "dry-run", "d", false, "Dry run. Outputs the generated yaml files without deploying")
	releaseCmd.Flags().StringVarP(&ns, "namespace", "n", "default", "Namespace to deploy to")
	releaseCmd.Flags().StringVarP(&packfile, "packfile", "p", "", "The values yaml file to use")
	releaseCmd.Flags().BoolVarP(&xdebug, "xdebug", "x", false, "Enables xdebug (for dev environments only)")
	releaseCmd.Flags().BoolVar(&noExecute, "no-execute", false, "Echoes helm upgrade command, but does not execute")
	releaseCmd.Flags().StringVar(&packageIDOverride, "packageid", "", "Package ID. Overrides default set based on environment")
	releaseCmd.Flags().StringVar(&optSetValues, "set", "", "(From Helm) set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
}

func getK8sCurrContext() ([]byte, error) {
	cmdName := "kubectl"
	cmdArgs := []string{"config", "current-context"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	check(err)
	return cmdOut, err
}

func getXDebugHost() ([]byte, error) {
	cmdName := "ipconfig"
	cmdArgs := []string{"getifaddr", "en0"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	check(err)
	return cmdOut, err
}

func useK8sCurrContext(context string) ([]byte, error) {
	cmdName := "kubectl"
	cmdArgs := []string{"config", "use-context", context}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	check(err)
	return cmdOut, err
}

func execHelmUpgradeCmd(fullReleaseName string, appPath string, setValues string, packfileFullPath string, packfile string, ns string) {
	msg := "Running helm upgrade"

	fullSetValues := setValues
	if len(optSetValues) > 0 {
		fullSetValues += "," + optSetValues
	}

	releasePath := viper.GetString("ReleasePath")
	globalPath := releasePath + "/.global/"
	globalValuesPath := globalPath + "values.yaml"
	globalValuesEnvPath := globalPath + packfile

	var fullPackFiles string

	//precedence should go like this:
	//values.env.yaml, .global/values.yaml, .global/values.env.yaml
	//right values files override left
	if pathExists(packfileFullPath) {
		fullPackFiles = packfileFullPath
	} else {
		echoWarningMessage(packfile + " does not exist. Running helm upgrade with values.yaml only\n")
	}

	if pathExists(globalPath) {

		if len(fullPackFiles) > 0 {
			fullPackFiles += ","
		}

		fullPackFiles += globalValuesPath

		if pathExists(globalValuesEnvPath) {
			fullPackFiles += "," + globalValuesEnvPath
		}

	}

	cmdName := "helm"
	cmdArgs := []string{
		"upgrade", fullReleaseName, "--install", appPath, "--set", fullSetValues, "--namespace", ns}

	if dryrun {
		cmdArgs = append(cmdArgs, "--dry-run", "--debug")
		msg += " (dry run)"
	}

	if len(fullPackFiles) > 0 {
		cmdArgs = append(cmdArgs, "--values", fullPackFiles)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdString := strings.Join(cmd.Args, " ")
	echoGoodMessage(cmdString)
	confirm := true

	if !noExecute {
		if env == "production" && !dryrun {
			confirm = askForConfirmation(fullReleaseName)
		}
		if confirm {
			fmt.Printf("\n%s\n", msg)
			out, _ := cmd.CombinedOutput()
			fmt.Printf("%s", out)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pathExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)

}

func echoWarningMessage(msg string) {
	fmt.Printf("%s%s%s", colorYellow, msg, colorNone)
}

func echoGoodMessage(msg string) {
	fmt.Printf("%s%s%s", colorGreen, msg, colorNone)
}

func askForConfirmation(appName string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("")
		echoWarningMessage("Do you really want to deploy `" + appName + "` to production? [y/n]: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
