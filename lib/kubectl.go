package cmd

import (
	"io"
	"os/exec"
)

var cmName = "boatswain-config"
var cmColumnPath = ".data.config"

type Kubectl struct{}

func (k *Kubectl) GetConfigMap() []byte {

	name := "kubectl"
	cmdArgs := []string{"get", "secret", cmName, "-o", "custom-columns=Config:" + cmColumnPath}
	cmd := exec.Command(name, cmdArgs...)

	out, err := cmd.CombinedOutput()
	k.checkError(err)

	return out
}

func (k *Kubectl) UpdateConfigMap(manifest string) []byte {
	name := "kubectl"
	cmdArgs := []string{"apply", "-f", "-"}
	cmd := exec.Command(name, cmdArgs...)
	stdin, _ := cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, manifest)
	}()

	out, err := cmd.CombinedOutput()
	k.checkError(err)

	return out
}

func (k *Kubectl) checkError(err error) {
	if err != nil {
		panic(err)
	}
}
