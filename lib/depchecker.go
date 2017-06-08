package lib

import (
	"os"
	"os/exec"

	logging "github.com/op/go-logging"
)

type DepChecker struct {
	log logging.Logger
}

func NewDepChecker(logger logging.Logger) *DepChecker {
	depChecker := DepChecker{log: logger}
	return &depChecker
}

func (dep *DepChecker) CheckDepDocker() {
	cmdDocker := "docker"
	cmdDockerArgs := []string{"info"}
	_, err := exec.Command(cmdDocker, cmdDockerArgs...).CombinedOutput()
	if err != nil {
		dep.log.Critical("Missing dependency: Docker is not running. Please start docker and try again. To install Docker for Mac, go here: https://docs.docker.com/docker-for-mac/install/")
		os.Exit(1)
	}
}

func (dep *DepChecker) CheckDepAWS() {
	dep.CheckInPath("aws", "http://docs.aws.amazon.com/cli/latest/userguide/installing.html")
}

func (dep *DepChecker) CheckDepHelm() {
	dep.CheckInPath("helm", "https://github.com/kubernetes/helm/releases/latest")
}

func (dep *DepChecker) CheckDepKubectl() {
	dep.CheckInPath("kubectl", "https://kubernetes.io/docs/tasks/tools/install-kubectl")
}

func (dep *DepChecker) CheckInPath(name string, installUrl string) {
	cmd := "which"
	cmdArgs := []string{name}
	out, err := exec.Command(cmd, cmdArgs...).CombinedOutput()
	if err != nil {
		dep.log.Criticalf("Missing dependency: %s is not in installed or is not in your PATH. To install %s, go here: %s", name, name, installUrl)
		os.Exit(1)
	} else {
		dep.log.Debugf("Found dependency %s in %s", name, out)
	}
}
