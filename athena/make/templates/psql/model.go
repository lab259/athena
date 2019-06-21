package templates_psql

import "github.com/lab259/athena/athena/util/template"

type ModelTemplateData struct {
	Project  string
	Model    string
	Table    string
	Fields   []string
	WithCRUD bool
}

var ModelTemplate = template.New("model.go", `package models

import (
	kallax "gopkg.in/src-d/go-kallax.v1"
)

// {{.Model}} TODO
type {{.Model}} struct {
	kallax.Model         `+"`"+`table:"{{.Table}}"`+"`"+`
	kallax.Timestamps
	ID kallax.ULID  `+"`"+`json:"id" pk:""`+"`"+`
	{{range .Fields}}{{formatField .}}  `+"`"+`json:"{{formatFieldTag .}}"{{if hasValidation .}} validate:"{{formatValidation .}}"{{end}}`+"`"+`
	{{end}}
}
`)
