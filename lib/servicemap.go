package lib

import (
	"io/ioutil"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

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

func NewServiceMapConfig() *ServiceMapConfig {
	config := ServiceMapConfig{}
	return &config
}

func NewStagingServiceMapConfig() *ServiceMapConfig {
	config := NewServiceMapConfig()
	path := viper.GetString("release")
	fullPath := path + "/.servicemap/staging.yaml"
	valuesBytes, err := ioutil.ReadFile(fullPath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(valuesBytes, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func (s *ServiceMapConfig) GetServiceMap(name string) *ServiceMap {
	for _, smap := range s.ServiceMaps {
		if smap.Name == name {
			return &smap
		}
	}
	return nil
}

func (s *ServiceMap) GetEnvironmentVars(packageID string) map[string]string {
	env := make(map[string]string)

	for _, svc := range s.Test {
		env[svc] = packageID + "-" + svc
	}

	for _, svc := range s.Staging {
		env[svc] = "staging-" + svc
	}

	return env
}
