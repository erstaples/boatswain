package cmd

import (
	"bytes"
	"html/template"

	yaml "gopkg.in/yaml.v2"
)

type StagingConfigMap struct {
	Config []StagingConfigMapEntry `yaml:"Config"`
}

type StagingConfigMapEntry struct {
	Name                string   `yaml:"Name"`
	HelmDeployments     []string `yaml:"HelmDeployments"`
	Ingress             string   `yaml:"Ingress"`
	AutogenConfigs      []string `yaml:"AutogenConfigs"`
	CloudFormationStack string   `yaml:"CloudFormationStack"`
}

//RenderTemplate returns a ConfigMap object string
func (m *StagingConfigMap) RenderTemplate() string {
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
	  AutogenConfigs:
	    {{-range .AutogenConfigs }}
        - {{ . }}
		{{- end }}
	  Ingress: {{ .Ingress }}
	  CloudFormationStack: {{ .CloudFormationStack }}
  {{- end }}
`)

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, m)
	s := doc.String()

	if err != nil {
		panic(err)
	}
	return s
}

//LoadConfigMap loads config map from k8s and unmarshals data
func (m *StagingConfigMap) LoadConfigMap() {
	var k Kubectl
	out := k.GetConfigMap()
	yaml.Unmarshal(out, &m)
}

//AddConfig appends a new StagingConfigMapEntry object to config list
func (m *StagingConfigMap) AddConfig(c StagingConfigMapEntry) {
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
func (m *StagingConfigMap) Save() {
	var k Kubectl
	manifest := m.RenderTemplate()
	k.UpdateConfigMap(manifest)
}
