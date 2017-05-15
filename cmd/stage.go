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

type StagingValuesYAML struct {
	ImageTag             string
	ServiceMap           []string `yaml:"ServiceMap"`
	Env                  map[string]string
	CloudFormationValues map[string]string `yaml:"CloudFormationValues"`
}

type ServiceMapConfig struct {
	ServiceMaps []ServiceMap      `yaml:"ServiceMaps"`
	Ingress     ServiceMapIngress `yaml:"Ingress"`
}

type ServiceMap struct {
	Name                   string   `yaml:"Name"`
	Test                   []string `yaml:"Test"`
	Staging                []string `yaml:"Staging"`
	CloudFormationTemplate string   `yaml:"CloudFormationTemplate"`
}

type ServiceMapIngress struct {
	Template string `yaml:"Template"`
	Service  string `yaml:"Service"`
	Port     string `yaml:"Port"`
}

// stageCmd represents the stage command
var StageCmd = &cobra.Command{
	Use:   "stage [push|delete]",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use boatswain stage [push|delete]")
	},
}

func init() {
	RootCmd.AddCommand(StageCmd)
}
