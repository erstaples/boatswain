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
	"fmt"
	"io/ioutil"

	"os/user"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type KubeConfig struct {
	ApiVersion     string    `yaml:"apiVersion"`
	Clusters       []Cluster `yaml:"clusters"`
	CurrentContext string    `yaml:"current-context"`
	Kind           string    `yaml:"kind"`
	Contexts       []Context `yaml:"contexts"`
	Users          []User    `yaml:"users"`
}

type Context struct {
	Context map[string]string `yaml:"context"`
	Name    string
}

type Cluster struct {
	Cluster map[string]string `yaml:"cluster"`
	Name    string            `yaml:"name"`
}

type User struct {
	User map[string]string `yaml:"user"`
	Name string            `yaml:"name"`
}

var out string

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge <contextName> <mergePath>",
	Short: "Merge two kubeconfig files",
	Long: `Merge two kubeconfig files. By default, merge from <mergePath> into ${HOME}/.kube/config. 
Sets the kubeconfig context based on <contextName>. For example:

Merge /Users/foo/prod into ${HOME}/.kube/config with context name "production":
boatswain kubeconfig merge production /Users/foo/prod

Merge /Users/foo/prod into ~/diff/kube/config:
boatswain kubeconfig merge production /Users/foo/prod --out ~/diff/kube/config`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("Invalid arguments. Use --help option for usage")
		}

		contextName := args[0]
		userName := contextName + "-user"
		mergePath := args[1]
		sourcePath := out

		var mergeConfig KubeConfig
		var sourceConfig KubeConfig
		mergeBytes, err := ioutil.ReadFile(mergePath)
		if err != nil {
			panic(err)
		}

		yaml.Unmarshal(mergeBytes, &mergeConfig)

		sourceBytes, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			panic(err)
		}
		yaml.Unmarshal(sourceBytes, &sourceConfig)

		mergeConfig.Contexts[0].Name = contextName
		mergeConfig.Contexts[0].Context["user"] = userName
		mergeConfig.Users[0].Name = userName

		sourceConfig.Clusters = append(sourceConfig.Clusters, mergeConfig.Clusters[0])
		sourceConfig.Contexts = append(sourceConfig.Contexts, mergeConfig.Contexts[0])
		sourceConfig.Users = append(sourceConfig.Users, mergeConfig.Users[0])

		sourceBytes, err = yaml.Marshal(sourceConfig)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile(sourcePath, sourceBytes, 0666)
	},
}

/*
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <auth-data>
    server: https://54.0.0.3:6443
  name: spcab12abc
contexts:
- context:
    cluster: spcab12abc
    user: stackpoint
  name: stackpoint
current-context: stackpoint
kind: Config
preferences: {}
users:
- name: stackpoint
  user:
    client-certificate-data: <cert-data>
    client-key-data: <key-data>
*/

func init() {
	KubeconfigCmd.AddCommand(mergeCmd)
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	defaultOut := usr.HomeDir + "/.kube/config"
	mergeCmd.Flags().StringVarP(&out, "out", "o", defaultOut, "Output kubeconfig file.")
}
