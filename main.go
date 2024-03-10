package main

import (
	_ "embed"
	"html/template"
	"log"
	"os"
	"strings"
)

//go:embed templates/template.svg
var svgTemplate string

func main() {
	data := struct {
		ImageURL  string
		TextLine1 string
		TextLine2 string
	}{
		ImageURL:  "./templates/image.png",
		TextLine1: "Text 1",
		TextLine2: "Text 2",
	}
	data.ImageURL = strings.Replace(data.ImageURL, "&", "&amp;", -1)

	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		return
	}

	file, err := os.Create("output.svg")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		return
	}
}
