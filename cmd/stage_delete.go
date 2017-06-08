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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/medbridge/boatswain/lib"
	"github.com/medbridge/boatswain/utilities"
	"github.com/spf13/cobra"
)

var stageDeleteCmd = &cobra.Command{
	Use:   "delete [appnames] [domain]",
	Short: "Delete a deployment from staging",
	Long: `Delete an application or bundle of applications from staging
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStageDelete(args)
	},
}

func init() {
	StageCmd.AddCommand(stageDeleteCmd)
}

//RunStageDelete performs cleanup on stagings. Deletes CloudFormation stack, Ingress, Helm deployment, and staging values files
func RunStageDelete(args []string) {
	var confMap lib.StagingConfigMap
	domain := args[0]
	entry := confMap.Find(domain)

	if entry == nil {
		Logger.Criticalf("Staging release %s not found", domain)
		return
	}

	//helm deletes
	for _, helm := range entry.HelmDeployments {
		fullName := domain + "-" + helm
		utilities.ExecStreamOut("helm", []string{"delete", fullName}, Logger, false)
	}

	//delete ing
	utilities.ExecStreamOut("kubectl", []string{"delete", "ing", entry.Ingress}, Logger, false)

	//delete values
	path := lib.GetReleasePath()
	for _, deploy := range entry.HelmDeployments {
		stagingFile := path + "/" + deploy + "/autogenerated/values.staging." + domain + ".yaml"
		err := os.Remove(stagingFile)
		if err != nil {
			Logger.Warningf("%s not found.", stagingFile)
			os.Exit(1)
		}
	}

	//delete stack
	if len(entry.CloudFormationStack) > 0 {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
		}))

		svc := cloudformation.New(sess)
		deleteStackParams := &cloudformation.DeleteStackInput{
			StackName: aws.String(entry.CloudFormationStack),
		}

		out, err := svc.DeleteStack(deleteStackParams)
		if err != nil {
			Logger.Errorf("%s", err)
			panic(err)
		} else {
			Logger.Infof("%s", out)
		}

	}

	confMap.Delete(domain)
}
