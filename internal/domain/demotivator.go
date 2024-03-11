package domain

import (
	"demotivator-generator/templates"
	"io"
)

type Demotivator struct {
	ImageBase64 string
	TextLine1   string
	TextLine2   string
}

func (d Demotivator) Generate(resultWriter io.Writer) error {
	template, err := templates.GetTemplate()
	if err != nil {
		return err
	}

	err = template.Execute(resultWriter, d)
	if err != nil {
		return err
	}
	return nil
}
