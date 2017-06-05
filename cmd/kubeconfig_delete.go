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

	"github.com/medbridge/boatswain/lib"
	"github.com/spf13/cobra"
)

// deleteCmd represents the merge command
var deleteCmd = &cobra.Command{
	Use: "delete <contextName>",
	Short: `
	Delete a context in a kubeconfig file`,
	Long: `

	Delete a context in a kubeconfig file. 
	By default, target kubeconfig file is ${HOME}/.kube/config. 
	For example:

	Delete context "staging" from ${HOME}/.kube/config:
	boatswain kubeconfig delete staging

	Delete context "staging" from ${HOME}/my/config:
	boatswain kubeconfig delete staging -f ${HOME}/my/config
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Invalid arguments. Use --help option for usage")
		}
		contextName := args[0]
		sourcePath := file
		sourceConfig := lib.NewKubeConfig(sourcePath)

		found := sourceConfig.DeleteContext(contextName)
		if found {
			sourceConfig.WriteFile()
			fmt.Printf("Context deleted: %s", contextName)
		} else {
			fmt.Printf("Context not found: %s", contextName)
		}
	},
}

func init() {
	KubeconfigCmd.AddCommand(deleteCmd)
}
