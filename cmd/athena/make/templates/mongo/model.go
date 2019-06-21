package templates_mongo

import "github.com/lab259/athena/cmd/athena/util/template"

type ModelTemplateData struct {
	Project        string
	Model          string
	Collection     string
	WithRepository bool
	WithCRUD       bool
	Fields         []string
}

var ModelTemplate = template.New("model.go", `package models

import (
	"context"
	{{if .WithRepository}}
	"github.com/lab259/repository"
	mgorscsrv "github.com/lab259/athena/rscsrv/mgo"
	"github.com/gofrs/uuid"
	{{end}}
)

// {{.Model}} TODO
type {{.Model}} struct {
	{{if .WithRepository}}ID uuid.UUID  `+"`"+`json:"id" bson:"_id"`+"`"+`{{end}}
	{{range .Fields}}{{formatField .}}  `+"`"+`json:"{{formatFieldTag .}}" bson:"{{formatFieldTag .}}"{{if hasValidation .}} validate:"{{formatValidation .}}"{{end}}`+"`"+`
	{{end}}
}

{{if .WithRepository}}
// New{{.Model}}Repository returns a Repository instance for {{.Model}} model
func New{{.Model}}Repository(ctx context.Context) *repository.Repository {
	return repository.NewRepository(repository.RepositoryConfig{
		Collection:  "{{.Collection}}",
  		QueryRunner: mgorscsrv.DefaultMgoService,
	})
}
{{end}}
`)
