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
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var filePath string

var genCredsCmd = &cobra.Command{
	Use:   "gencreds <accessKeyCsvFile>",
	Short: "Enable and configure awsecr-creds minikube addon",
	Long: `Example:
	
	boatswain gencreds ~/Desktop/accessKeys.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Required argument: filePath")
		}

		execEnableAddonCmd()
		execGenSecretCmd(args[0])
	},
}

func init() {
	RootCmd.AddCommand(genCredsCmd)
}

func execEnableAddonCmd() {
	cmdName := "minikube"
	cmdArgs := []string{"addons", "enable", "awsecr-creds"}
	out, _ := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	fmt.Printf("%s", out)
}

func execGenSecretCmd(filePath string) {
	creds, _ := parseCreds(filePath)
	awsAccount := "191682557156"
	cmdName := "kubectl"
	cmdArgs := []string{
		"--namespace=kube-system",
		"create", "secret", "generic",
		"awsecr-creds",
		"--from-literal=AWS_ACCESS_KEY_ID=" + creds[0],
		"--from-literal=AWS_SECRET_ACCESS_KEY=" + creds[1],
		"--from-literal=aws-account=" + awsAccount}
	out, _ := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	fmt.Printf("%s", out)

}

func parseCreds(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Read()             //dump header row
	row, err := csvr.Read() //extract the juicy bits
	return row, err
}
