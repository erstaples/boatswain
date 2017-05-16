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
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	cm "github.com/medbridge/boatswain/lib/configmap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stageDeleteCmd = &cobra.Command{
	Use:   "delete [appnames] [domain]",
	Short: "Delete an application(s) from staging",
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
	var confMap cm.StagingConfigMap
	domain := args[0]
	entry := confMap.Find(domain)

	if entry == nil {
		fmt.Println("\nStaging release not found")
		return
	}

	//helm deletes
	for _, helm := range entry.HelmDeployments {
		fullName := domain + "-" + helm
		out, err := exec.Command("helm", "delete", fullName).CombinedOutput()
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			fmt.Printf("%s", out)
		}
	}

	//delete ing
	out, err := exec.Command("kubectl", "delete", "ing", entry.Ingress).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("%s", out)
	}

	//delete values
	path := viper.GetString("release")
	for _, deploy := range entry.HelmDeployments {
		stagingFile := path + "/" + deploy + "/autogenerated/values.staging." + domain + ".yaml"
		err = os.Remove(stagingFile)
		if err != nil {
			panic(err)
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
			panic(err)
		} else {
			fmt.Printf("%s", out)
		}

	}

	confMap.Delete(domain)
}
