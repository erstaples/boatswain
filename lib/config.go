package lib

import (
	"io/ioutil"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ReleasePath string  `yaml:"ReleasePath"`
	Builds      []Build `yaml:"Builds"`
}

func GetReleasePath() string {
	return viper.GetString("ReleasePath")
}

func LoadConfig() *Config {
	config := Config{}
	configPath := viper.ConfigFileUsed()
	yamlFile, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(yamlFile, &config)
	return &config
}
