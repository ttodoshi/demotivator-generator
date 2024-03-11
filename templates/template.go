package templates

import (
	_ "embed"
	"html/template"
)

//go:embed template.svg
var svgTemplate string

func GetTemplate() (*template.Template, error) {
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
