package lib

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/*
sample kubeconfig format, for reference:

apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <auth-data>
    server: https://54.0.0.3:6443
  name: spcab12abc
contexts:
- context:
    cluster: spcab12abc
    user: stackpoint
  name: stackpoint
current-context: stackpoint
kind: Config
preferences: {}
users:
- name: stackpoint
  user:
    client-certificate-data: <cert-data>
    client-key-data: <key-data>
*/

type KubeConfig struct {
	ApiVersion     string    `yaml:"apiVersion"`
	Clusters       []Cluster `yaml:"clusters"`
	CurrentContext string    `yaml:"current-context"`
	Kind           string    `yaml:"kind"`
	Contexts       []Context `yaml:"contexts"`
	Users          []User    `yaml:"users"`
	Path           string
}

type Context struct {
	Context map[string]string `yaml:"context"`
	Name    string
}

type Cluster struct {
	Cluster map[string]string `yaml:"cluster"`
	Name    string            `yaml:"name"`
}

type User struct {
	User map[string]string `yaml:"user"`
	Name string            `yaml:"name"`
}

func NewKubeConfig(path string) *KubeConfig {
	kc := KubeConfig{Path: path}
	readInConfig(&kc)
	return &kc
}

func (kc *KubeConfig) WriteFile() {
	bytes, err := yaml.Marshal(kc)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(kc.Path, bytes, 0666)
}

func (kc *KubeConfig) DeleteContext(name string) bool {

	targetContext := Context{}
	var contextIndex int
	var clusterIndex int
	var userIndex int
	found := false
	for i, context := range kc.Contexts {
		if context.Name == name {
			targetContext = context
			contextIndex = i
			found = true
		}
	}
	if found == false {
		return found
	}

	for i, cluster := range kc.Clusters {
		if cluster.Name == targetContext.Context["cluster"] {
			clusterIndex = i
		}
	}

	for i, user := range kc.Users {
		if user.Name == targetContext.Context["user"] {
			userIndex = i
		}
	}

	kc.Contexts = append(kc.Contexts[:contextIndex], kc.Contexts[contextIndex+1:]...)
	kc.Users = append(kc.Users[:userIndex], kc.Users[userIndex+1:]...)
	kc.Clusters = append(kc.Clusters[:clusterIndex], kc.Clusters[clusterIndex+1:]...)

	return found
}

func (kc *KubeConfig) ListContexts() {
	for _, context := range kc.Contexts {
		fmt.Printf("\n- %s\n  user: %s\n  cluster: %s", context.Name, context.Context["user"], context.Context["cluster"])
	}
}

func (kc *KubeConfig) MergeContext(mergeConfig *KubeConfig, contextName string) {
	userName := contextName + "-user"
	mergeConfig.Contexts[0].Name = contextName
	mergeConfig.Contexts[0].Context["user"] = userName
	mergeConfig.Users[0].Name = userName

	kc.Clusters = append(kc.Clusters, mergeConfig.Clusters[0])
	kc.Contexts = append(kc.Contexts, mergeConfig.Contexts[0])
	kc.Users = append(kc.Users, mergeConfig.Users[0])
}

func readInConfig(kc *KubeConfig) {
	bytes, err := ioutil.ReadFile(kc.Path)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(bytes, &kc)
}
