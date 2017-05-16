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
	"github.com/medbridge/mocking/factories"
	"github.com/spf13/cobra"
)

var build = lib.Build{}
var packageID string

var stagePushCmd = &cobra.Command{
	Use:   "push [appnames] [domain]",
	Short: "Push an application(s) to staging",
	Long: `Push an application or bundle of applications to staging

	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStagePush(args)
	},
}

func init() {
	StageCmd.AddCommand(stagePushCmd)
}

func RunStagePush(args []string) {
	if len(args) != 2 {
		fmt.Printf("Unexpected number of args. Expected 2, got %c", len(args))
		return
	}

	serviceMapName := args[0]
	packageID = args[1]
	smapConfig := lib.NewStagingServiceMapConfig()
	smap := smapConfig.GetServiceMap(serviceMapName)
	builds := lib.GetBuilds(*smap)
	cloudformation := lib.CloudFormationTemplate{Output: make(map[string]string)}
	env := smap.GetEnvironmentVars(packageID)
	imageTags := make(map[string]string)

	if len(smap.CloudFormationTemplate) > 0 {
		cloudformation = *lib.NewCloudFormationTemplate(smap.CloudFormationTemplate)
		cloudformation.CreateStack(packageID)
	}

	for _, build := range builds {
		imageTags[build.Name] = build.Exec()
	}

	helmDeploys := []string{}
	for _, svc := range smap.Test {
		values := lib.NewValues(packageID, svc, imageTags[svc], env)
		values.CloudFormationValues = cloudformation.Output
		runRelease(svc, values.Write())
		helmDeploys = append(helmDeploys, svc)
	}

	genIngress(*smapConfig)

	lib.NewStagingConfigMap().AddConfig(
		lib.StagingConfigMapEntry{
			CloudFormationStack: cloudformation.StackName,
			HelmDeployments:     helmDeploys,
			Name:                packageID,
			Ingress:             smapConfig.Ingress.RenderHostName(packageID),
		})

}

func runRelease(name string, valuesFile string) {
	args := []string{name}
	options := ReleaseOptions{
		Environment:       "staging",
		DryRun:            false,
		Namespace:         "default",
		Packfile:          valuesFile,
		Xdebug:            false,
		NoExecute:         false,
		PackageIDOverride: packageID,
	}

	RunRelease(args, options)
}

func genIngress(config lib.ServiceMapConfig) {
	cmdFactory := &factories.CommandFactory{}
	args := []string{config.Ingress.RenderHostName(packageID)}
	options := GenIngressFlags{
		Service:     packageID + "-" + config.Ingress.Service,
		EnableTLS:   false,
		ServicePort: config.Ingress.Port,
	}

	RunGenIngress(args, cmdFactory, options)
}
