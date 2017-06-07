package lib

import (
	"bytes"
	"html/template"
	"strings"

	logging "github.com/op/go-logging"

	yaml "gopkg.in/yaml.v2"
)

var logger logging.Logger

type StagingConfigMap struct {
	Config   []StagingConfigMapEntry `yaml:"Config"`
	IsLoaded bool
}

type StagingConfigMapEntry struct {
	Name                string   `yaml:"Name"`
	HelmDeployments     []string `yaml:"HelmDeployments"`
	Ingress             string   `yaml:"Ingress"`
	CloudFormationStack string   `yaml:"CloudFormationStack"`
}

func NewStagingConfigMap() *StagingConfigMap {
	m := StagingConfigMap{}
	m.LoadConfigMap()
	return &m
}

//RenderTemplate returns a ConfigMap object string
func (m *StagingConfigMap) RenderTemplate() string {
	/**
	  This template is fairly picky in terms of indentations. If you need to
	  edit this, make sure it outputs a properly-formatted yaml string. And make sure
	  new columns are indented with spaces, not tabs
	  **/
	tmpl, err := template.New("configmap").Parse(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: boatswain-config
data:
  config: |
  {{- range .Config }}
    - Name: {{ .Name }}
      HelmDeployments:
	    {{- range .HelmDeployments }}
        - {{ . }}
			{{- end }}
      Ingress: {{ .Ingress }}
      CloudFormationStack: {{ .CloudFormationStack }}
  {{- end }}
`)

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, m)
	s := doc.String()
	logger.Debugf("%s", s)

	if err != nil {
		panic(err)
	}
	return s
}

//LoadConfigMap loads config map from k8s and unmarshals data
func (m *StagingConfigMap) LoadConfigMap() {
	var k Kubectl
	out := k.GetConfigMap()
	//the column header won't have the colon, which we need to unmarshal, so add it here
	out = []byte(strings.Replace(string(out), "Config", "Config:", 1))
	yaml.Unmarshal(out, &m)
	m.IsLoaded = true
}

//AddConfig appends a new StagingConfigMapEntry object to config list. If the entry with the same name exists, it replaces the existing entry
func (m *StagingConfigMap) AddConfig(c StagingConfigMapEntry, log logging.Logger) {
	logger = log
	found := false
	for i, entry := range m.Config {
		if entry.Name == c.Name {
			m.Config[i] = c
			found = true
		}
	}
	if found == false {
		m.Config = append(m.Config, c)
	}
	m.Save()
}

//Save renders ConfigMap template and pushes it to k8s
func (m *StagingConfigMap) Save() string {
	var k Kubectl
	manifest := m.RenderTemplate()
	k.UpdateConfigMap(manifest)
	return manifest
}

func (m *StagingConfigMap) Find(name string) *StagingConfigMapEntry {
	if m.IsLoaded == false {
		m.LoadConfigMap()
	}

	for _, entry := range m.Config {
		if entry.Name == name {
			return &entry
		}
	}

	return nil
}

func (m *StagingConfigMap) Delete(name string) {
	if m.IsLoaded == false {
		m.LoadConfigMap()
	}

	for i, entry := range m.Config {
		if entry.Name == name {
			m.Config = append(m.Config[:i], m.Config[i+1:]...)
		}
	}

	m.Save()
}
