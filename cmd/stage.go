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
	"github.com/spf13/cobra"
)

// stageCmd represents the stage command
var stageCmd = &cobra.Command{
	Use:   "stage [push|delete] [appnames] [domainname]",
	Short: "A brief description of your command",
	Long:  ``,
}

var stagePushCmd = &cobra.Command{
	Use:   "stage push [appnames] [domainname]",
	Short: "Push an application(s) to staging",
	Long: `Push an application or bundle of applications to staging

	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStagePush(args)
	},
}

var stageDeleteCmd = &cobra.Command{
	Use:   "stage delete [appnames] [domainname]",
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stageCmd.Flags().StringVar(&StageFlags."toggle", "t", false, "Help message for toggle")

}

func RunStagePush(args []string) {

}

func RunStageDelete(args []string) {

}
