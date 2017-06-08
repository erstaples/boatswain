package lib

import (
	"os"
	"os/exec"

	"github.com/medbridge/boatswain/utilities"
	logging "github.com/op/go-logging"
)

type Build struct {
	Name     string `yaml:"Name"`
	Path     string `yaml:"Path"`
	Rootpath string `yaml:"RootPath"`
	ImageTag string
	Logger   *logging.Logger
}

//Exec Runs the build shell script at Path
func (b *Build) Exec() string {
	cmdName := "/bin/bash"
	cmdArgs := []string{b.Path, "push"}

	os.Chdir(b.Rootpath)
	b.Logger.Infof("Running %s", append([]string{cmdName}, cmdArgs...))
	utilities.ExecStreamOut(cmdName, cmdArgs, *b.Logger, true)

	b.SetImageTag()
	return b.ImageTag
}

//SetImageTag gets a git commit sha from RootPath and sets the ImageTag
func (b *Build) SetImageTag() {
	os.Chdir(b.Rootpath)
	cmdName := "git"
	cmdArgs := []string{"show", "-s", "--pretty=format:%h"}
	out, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		panic(err)
	}
	cmd := append([]string{cmdName}, cmdArgs...)
	b.Logger.Debugf("%s: %s", cmd, out)
	b.ImageTag = string(out[:])
}

//GetBuilds iterates over a service map and returns an array of Build objects that are needed for the release. Builds are defined in the ~/.boatswain.yaml config file
func GetBuilds(smap ServiceMap, logger logging.Logger) []Build {
	config := LoadConfig(&logger)
	var builds []Build
	for _, svc := range smap.Test {
		for _, build := range config.Builds {
			if svc == build.Name {
				builds = append(builds, build)
			}
		}
	}
	return builds
}
