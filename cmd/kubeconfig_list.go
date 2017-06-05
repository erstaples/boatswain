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
	"github.com/medbridge/boatswain/lib"
	"github.com/spf13/cobra"
)

// listCmd represents the merge command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available contexts in a kubeconfig file",
	Long: `List available contexts in a kubeconfig file. By default, target kubeconfig file is ${HOME}/.kube/config. 
For example:

List contexts from ${HOME}/.kube/config:
boatswain kubeconfig list

List contexts from ${HOME}/my/config:
boatswain kubeconfig list -f ${HOME}/my/config`,
	Run: func(cmd *cobra.Command, args []string) {

		sourcePath := file

		sourceConfig := lib.NewKubeConfig(sourcePath)
		sourceConfig.ListContexts()
	},
}

func init() {
	KubeconfigCmd.AddCommand(listCmd)
}
