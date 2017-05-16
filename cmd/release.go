// Copyright © 2017 NAME HERE eric@medbridgeed.com
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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/medbridge/boatswain/lib"
	utils "github.com/medbridge/boatswain/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var env string
var dryrun bool
var ns string
var packfile string
var xdebug bool
var noExecute bool
var packageIDOverride string
var optSetValues string

type ReleaseOptions struct {
	Environment       string
	DryRun            bool
	Namespace         string
	Packfile          string
	Xdebug            bool
	NoExecute         bool
	PackageIDOverride string
	OptSetValues      string
}

type Config struct {
	Release string
	Builds  []lib.Build
}

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release appname",
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
		options := ReleaseOptions{
			Environment:       env,
			DryRun:            dryrun,
			Namespace:         ns,
			Packfile:          packfile,
			Xdebug:            xdebug,
			NoExecute:         noExecute,
			PackageIDOverride: packageIDOverride,
			OptSetValues:      optSetValues,
		}
		RunRelease(args, options)
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

func RunRelease(args []string, options ReleaseOptions) {
	if len(args) != 1 {
		fmt.Println("Invalid arguments.")
		return
	}

	releaseName := args[0]

	var xdebugHost string
	var fullReleaseName string
	var packageId string
	var envPackfile string

	environments := []string{"development", "staging", "production"}

	// based on environment, set default packageId, packfile, and context
	switch options.Environment {
	case environments[0]:
		if len(packageId) == 0 {
			packageId = "dev"
		}
		envPackfile = "values.env.yaml"
		useK8sCurrContext("minikube")
	case environments[1]:
		if len(packageId) == 0 {
			packageId = "staging"
		}
		envPackfile = "values.staging.yaml"
		useK8sCurrContext("staging")
	case environments[2]:
		if len(packageId) == 0 {
			packageId = "prod"
		}
		envPackfile = "values.prod.yaml"
		useK8sCurrContext("production")
	default:
		fmt.Println("Invalid environment: " + env)
		os.Exit(1)
	}

	if len(options.PackageIDOverride) > 0 {
		packageId = options.PackageIDOverride
	}

	fullReleaseName = packageId + "-" + releaseName

	//xdebug option turned on, so get the xdebug host ip address
	if options.Xdebug {

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

	releasePath := viper.GetString("release")
	appPath := releasePath + "/" + releaseName

	fmt.Printf("boatswain | Deploying: %s\n", appPath)

	//build helm cmd
	//add standard values to be made available in the helm releases
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	setValues := "environment=" + options.Environment + ",packageId=" + packageId + ",timestamp=" + timestamp
	if options.Xdebug {
		setValues += ",xdebugHost=" + xdebugHost
	}

	//fully qualified path
	packfileFullPath := appPath + "/" + envPackfile

	execHelmUpgradeCmd(fullReleaseName, appPath, setValues, packfileFullPath, envPackfile, options)
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

func execHelmUpgradeCmd(fullReleaseName string, appPath string, setValues string, packfileFullPath string, envPackfile string, options ReleaseOptions) {
	msg := "Running helm upgrade"

	fullSetValues := setValues
	if len(options.OptSetValues) > 0 {
		fullSetValues += "," + options.OptSetValues
	}

	releasePath := viper.GetString("release")
	globalPath := releasePath + "/.global/"
	globalValuesPath := globalPath + "values.yaml"
	globalValuesEnvPath := globalPath + envPackfile

	var fullPackFiles string

	//precedence should go like this:
	//values.env.yaml, .global/values.yaml, .global/values.env.yaml, autogenerated/values.staging.foo.yaml
	//right values files override left
	if pathExists(packfileFullPath) {
		fullPackFiles = packfileFullPath
	} else {
		utils.EchoWarningMessage(envPackfile + " does not exist. Running helm upgrade with values.yaml only\n")
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

	if pathExists(options.Packfile) {
		fullPackFiles += "," + options.Packfile
	}

	cmdName := "helm"
	cmdArgs := []string{
		"upgrade", fullReleaseName, "--install", appPath, "--set", fullSetValues, "--namespace", options.Namespace}

	if options.DryRun {
		cmdArgs = append(cmdArgs, "--dry-run", "--debug")
		msg += " (dry run)"
	}

	if len(fullPackFiles) > 0 {
		cmdArgs = append(cmdArgs, "--values", fullPackFiles)
	}

	confirm := true
	cmdString := strings.Join(cmdArgs, " ")
	fmt.Printf("boatswain | %s", cmdString)

	if !options.NoExecute {
		if env == "production" && !dryrun {
			confirm = askForReleaseConfirmation(fullReleaseName)
		}
		if confirm {
			fmt.Printf("\n%s\n", msg)
			tries := 0
			for tries < 3 {
				tries++
				//need to retry because it's failing randomly... not ideal and we should remove this loop
				//the random failure is found
				cmd := exec.Command(cmdName, cmdArgs...)
				out, err := executeReleaseCmd(cmd)
				if err != nil && tries == 3 {
					fmt.Printf("%s", err)
				}

				if err == nil {
					fmt.Printf("%s", out)
					tries = 3
				}

			}
		}
	}
}

func executeReleaseCmd(cmd *exec.Cmd) ([]byte, error) {
	out, err := cmd.CombinedOutput()

	return out, err
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

func askForReleaseConfirmation(appName string) bool {
	msg := "Do you really want to deploy '" + appName + "' to production? [y/n]: "
	return utils.AskForConfirmation(msg)
}
