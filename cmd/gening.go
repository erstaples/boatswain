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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

var service string
var servicePort string
var enableTLS bool

type Ingress struct {
	HostName    string
	ServiceName string
	ServicePort string
	SecretName  string
}

// geningCmd represents the gening command
var geningCmd = &cobra.Command{
	Use:   "gening <hostname>",
	Short: "Generate ingress",
	Long: `Example:
	
	boatswain gening localmed.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if len(args) == 0 {
			fmt.Println("Missing argument: host")
		}
		host := args[0]
		var ingress Ingress

		if enableTLS {
			secretName := "tls-" + host
			fmt.Println(secretName)

			cmdName := "openssl"
			cmdArgs := []string{
				"req", "-x509",
				"-sha256", "-nodes", "-newkey", "rsa:4096",
				"-keyout", "tls.key",
				"-out", "tls.crt",
				"-days", "365",
				"-subj", "/CN=" + host}

			out, _ := exec.Command(cmdName, cmdArgs...).CombinedOutput()

			fmt.Printf("%s", out)

			cmdName = "kubectl"
			cmdArgs = []string{
				"create", "secret", "tls", secretName, "--cert=./tls.crt", "--key=./tls.key"}
			out, _ = exec.Command(cmdName, cmdArgs...).CombinedOutput()

			fmt.Printf("%s", out)
			os.Remove("tls.crt")
			os.Remove("tls.key")

			ingress = Ingress{HostName: host, ServiceName: service, ServicePort: servicePort, SecretName: secretName}
		} else {
			ingress = Ingress{HostName: host, ServiceName: service, ServicePort: servicePort}
		}

		var k8smanifest string
		if enableTLS {
			k8smanifest, _ = getHTTPSIngress(ingress)
		} else {
			k8smanifest, _ = getHTTPIngress(ingress)
		}

		ingCmd := exec.Command("kubectl", "apply", "-f", "-")
		stdin, _ := ingCmd.StdinPipe()

		fmt.Printf("%s", k8smanifest)
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, k8smanifest)

		}()

		out, _ := ingCmd.CombinedOutput()
		fmt.Printf("%s", out)
	},
}

func init() {
	RootCmd.AddCommand(geningCmd)

	geningCmd.Flags().StringVarP(&service, "service", "s", "dev-medbridge", "Service name")
	geningCmd.Flags().StringVarP(&servicePort, "port", "p", "80", "Service port")
	geningCmd.Flags().BoolVarP(&enableTLS, "tls", "t", false, "Enable TLS")

}

func getHTTPIngress(ing Ingress) (string, error) {
	tmpl, err := template.New("ingress-notls").Parse(`
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{.HostName}}
spec:
  rules:
  - host: {{.HostName}}
    http:
      paths:
      - backend:
          serviceName: {{.ServiceName}}
          servicePort: {{.ServicePort}}
`)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, ing)
	s := doc.String()
	return s, err
}

func getHTTPSIngress(ing Ingress) (string, error) {
	tmpl, err := template.New("ingress-tls").Parse(`
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{.HostName}}
spec:
  rules:
  - host: {{.HostName}}
    http:
      paths:
      - backend:
          serviceName: {{.ServiceName}}
          servicePort: {{.ServicePort}}
  tls:
  - hosts:
    - {{.HostName}}
    secretName: {{.SecretName}}
`)
	var doc bytes.Buffer
	tmpl.Execute(&doc, ing)
	err = tmpl.Execute(&doc, ing)
	s := doc.String()
	return s, err
}

/**

MATCHING_ROW=$(cat /etc/hosts | grep "^[0-9\.]*\s*${HOST_NAME}$")
HAS_ROW=$?
if [ $HAS_ROW -eq 1 ] && [ "$(kubectl config current-context)" == "minikube" ]; then
echo "$(minikube ip)  ${HOST_NAME}" | sudo tee -a /etc/hosts
echo "Your /etc/hosts file has been updated"
fi

**/
