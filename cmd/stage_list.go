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

	"github.com/medbridge/boatswain/lib"
	"github.com/spf13/cobra"
)

var stageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active stagings",

	Run: func(cmd *cobra.Command, args []string) {
		RunStageList(args)
	},
}

func init() {
	StageCmd.AddCommand(stageListCmd)
}

//RunStageList output all active stagings
func RunStageList(args []string) {
	var confMap lib.StagingConfigMap
	confMap.LoadConfigMap()
	for _, staging := range confMap.Config {
		fmt.Printf("\n* %s", staging.Name)
	}

}
