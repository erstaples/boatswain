package lib

import (
	"os"
	"os/exec"
	"strings"

	"github.com/medbridge/boatswain/utilities"
)

type Build struct {
	Name     string
	Path     string
	Targets  string
	Rootpath string
	ImageTag string
}

func (b *Build) SetImageTag() {
	os.Chdir(b.Rootpath)
	cmdName := "git"
	cmdArgs := []string{"show", "-s", "--pretty=format:%h"}
	out, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		panic(err)
	}
	b.ImageTag = string(out[:])
}

func (b *Build) RunBuild() string {
	cmdName := "/bin/bash"
	cmdArgs := []string{b.Path, "push"}

	targetsString := b.Targets
	targets := strings.Split(targetsString, ",")

	if len(targets) > 0 {
		for _, target := range targets {
			cmdArgs = append(cmdArgs, target)
		}
	}
	os.Chdir(b.Rootpath)
	utilities.ExecStreamOut(cmdName, cmdArgs, "build.sh")

	b.SetImageTag()
	return b.ImageTag
}
