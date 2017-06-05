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

var overwrite bool

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge <contextName> <mergePath>",
	Short: "Merge two kubeconfig files",
	Long: `
	
Merge two kubeconfig files. 
By default, merge from <mergePath> into ${HOME}/.kube/config. 
Sets the kubeconfig context based on <contextName>. 

Important: merge only hanldes config in <mergePath> 
if it has a single context in it. 

Examples:

Merge /Users/foo/prod into ${HOME}/.kube/config with context name "production":
boatswain kubeconfig merge production /Users/foo/prod

Merge /Users/foo/prod into ~/diff/kube/config:
boatswain kubeconfig merge production /Users/foo/prod -f ~/diff/kube/config`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			fmt.Println("Invalid arguments. Use --help option for usage")
			return
		}

		contextName := args[0]
		mergePath := args[1]
		sourcePath := file

		mergeConfig := lib.NewKubeConfig(mergePath)
		sourceConfig := lib.NewKubeConfig(sourcePath)

		if !sourceConfig.ContextExists(contextName) || overwrite {
			sourceConfig.MergeContext(mergeConfig, contextName, overwrite)
			sourceConfig.WriteFile()
		} else {
			fmt.Printf("Context already exists. Delete context %s first, or use the overwrite flag (-o)", contextName)
		}

	},
}

func init() {
	KubeconfigCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing context if it already exists in the config file")
}
