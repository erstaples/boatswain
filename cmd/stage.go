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
	"fmt"

	"github.com/spf13/cobra"
)

// stageCmd represents the stage command
var StageCmd = &cobra.Command{
	Use:   "stage",
	Short: "Manage staging deployments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run `boatswain stage --help` for list of available subcommands.")
	},
}

func init() {
	RootCmd.AddCommand(StageCmd)
}
