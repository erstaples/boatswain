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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/medbridge/boatswain/lib"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugFormat = logging.MustStringFormatter(
	`%{color}%{shortfunc} ▶ %{level:.8s} %{color:reset} %{message}`,
)

var stdFormat = logging.MustStringFormatter(
	`%{color}%{level:.8s} %{color:reset} %{message}`,
)

// Version represents the app version. Used in `boatswain version` command
var Version = "v1.0.1-beta.2"

// Verbose output switch
var verbosity string
var Logger logging.Logger
var Deps *lib.DepChecker
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "boatswain",
	Short: "Utility to deploy applications and services",
	Long: `Provides a set of tools to deploy applications and services to a Kubernetes cluster. 

Boatswain does things like: 
* generates ingresses on the fly
* releases helm packages bootstrapped with the right environment and context
* provisions new AWS Kubernetes clusters with network policies and log aggregation


It makes some assumptions about your environment, which might be refactored in the near future
for portability:
* You have helm and kubectl installed and in your path
* You have a boatswain values repository
`, PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Logger = *logging.MustGetLogger("boatswain")
		backend := logging.NewLogBackend(os.Stderr, "boatswain: ", 0)
		format := stdFormat

		var logLevel logging.Level
		switch verbosity {
		case "debug":
			logLevel = logging.DEBUG
			format = debugFormat
			break
		case "info":
			logLevel = logging.INFO
			break
		case "critical":
			logLevel = logging.CRITICAL
			break
		default:
			logLevel = logging.INFO
		}
		backendFormatted := logging.NewBackendFormatter(backend, format)
		backendLeveled := logging.AddModuleLevel(backendFormatted)
		backendLeveled.SetLevel(logLevel, "")
		logging.SetBackend(backendLeveled)

		Deps = lib.NewDepChecker(Logger)
		Deps.CheckDepHelm()
		Deps.CheckDepKubectl()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "info", "Available values: debug,info,critical")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	usr, err := user.Current()
	homeDir := usr.HomeDir
	cfgFile = homeDir + "/.boatswain.yaml"
	viper.SetConfigFile(cfgFile)
	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
	} else {
		//if not found, initialize it
		genConfig()
	}
}
func genConfig() {
	fmt.Print("Enter path to boatswain/deployment folder (absolute path)\n")
	reader := bufio.NewReader(os.Stdin)
	path, _ := reader.ReadString('\n')

	yaml := "ReleasePath: " + path
	config := []byte(yaml)
	Logger.Infof("Creating config file at %s", cfgFile)
	err := ioutil.WriteFile(cfgFile, config, 0644)
	if err != nil {
		panic(err)
	}
	initConfig()
}
