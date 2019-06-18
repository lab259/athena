package template

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

func Write(tmpl *template.Template, data interface{}, name string) error {
	content := bytes.NewBuffer([]byte{})
	err := tmpl.Execute(content, data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, content.Bytes(), 0644)
}
