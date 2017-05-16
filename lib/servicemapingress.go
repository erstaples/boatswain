package lib

import (
	"bytes"
	"html/template"
)

type ServiceMapIngress struct {
	Template string `yaml:"Template"`
	Service  string `yaml:"Service"`
	Port     string `yaml:"Port"`
}

func (i *ServiceMapIngress) RenderHostName(packageID string) string {
	tmpl, _ := template.New("ingress").Parse(i.Template)
	var doc bytes.Buffer
	ingressName := struct {
		PackageID string
	}{packageID}
	err := tmpl.Execute(&doc, ingressName)
	if err != nil {
		panic(err)
	}
	return doc.String()
}
