package template

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"text/template"
)

func Write(tmpl *template.Template, data interface{}, name string) error {
	var content bytes.Buffer
	err := tmpl.Execute(&content, data)
	if err != nil {
		return err
	}
	formatted, err := format.Source(content.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, formatted, 0644)
}
