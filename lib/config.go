package lib

import (
	"io/ioutil"

	logging "github.com/op/go-logging"
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

func LoadConfig(logger *logging.Logger) *Config {
	config := Config{}
	configPath := viper.ConfigFileUsed()
	yamlFile, _ := ioutil.ReadFile(configPath)
	logger.Debugf("Config:\n%s", yamlFile)
	yaml.Unmarshal(yamlFile, &config)
	return &config
}
