package make

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/jinzhu/inflection"
	"github.com/lab259/athena/athena/util"
)

func Model(cmd *cli.Cmd) {
	var data modelTemplateData

	cmd.Spec = "[OPTIONS] MODEL FIELD..."

	cmd.StringArgPtr(&data.Model, "MODEL", "", "model name")
	cmd.StringsArgPtr(&data.Fields, "FIELD", []string{}, "field (eg.: Name:string)")
	cmd.BoolOptPtr(&data.WithRepository, "r repository", false, "with repository")

	cmd.Action = func() {
		data.Model = strcase.ToCamel(data.Model)
		data.Collection = strcase.ToSnake(inflection.Plural(data.Model))

		dir, err := os.Getwd()
		util.HandleError(err, "Unable to get current directory.")

		modelsDir := path.Join(dir, "models")

		err = os.MkdirAll(modelsDir, os.ModePerm)
		util.HandleError(err, "Unable to create package directory.")

		modelFile := fmt.Sprintf("%s.go", path.Join(modelsDir, strcase.ToSnake(data.Model)))

		content := bytes.NewBuffer([]byte{})
		err = mongoModelTemplate.Execute(content, &data)
		util.HandleError(err, "Unable to execute model template.")

		err = ioutil.WriteFile(modelFile, content.Bytes(), 0644)
		util.HandleError(err, "Unable to create model template.")

		fmt.Printf("%s was created.\n", modelFile)
	}
}

type modelTemplateData struct {
	Model          string
	Collection     string
	WithRepository bool
	Fields         []string
}

var mongoModelTemplate = template.Must(template.New("make:model:mongo").Funcs(template.FuncMap{
	"formatField": func(field string) (string, error) {
		return strings.Replace(field, ":", " ", 1), nil
	},
	"formatFieldTag": func(field string) (string, error) {
		values := strings.SplitN(field, ":", 2)
		if len(values) != 2 {
			return "", fmt.Errorf("invalid format for field: %s", field)
		}
		return strings.ReplaceAll(strcase.ToLowerCamel(values[0]), "ID", "Id"), nil
	},
}).Parse(`package models

import (
	"context"
	{{if .WithRepository}}
	"github.com/lab259/repository"
	"github.com/lab259/athena/rscsrv/mgo"
	"github.com/gofrs/uuid"
	{{end}}
)

// {{.Model}} TODO
type {{.Model}} struct {
	{{if .WithRepository}}ID uuid.UUID  ` + "`" + `json:"id" bson:"_id"` + "`" + `{{end}}
	{{range .Fields}}{{formatField .}}  ` + "`" + `json:"{{formatFieldTag .}}" bson:"{{formatFieldTag .}}"` + "`" + `
	{{end}}
}

{{if .WithRepository}}
// New{{.Model}}Repository returns a Repository instance for {{.Model}} model
func New{{.Model}}Repository(ctx context.Context) *repository.Repository {
	return repository.NewRepository({
		Collection:  "{{.Collection}}",
  		QueryRunner: &mgo.DefaultMgoService,
	})
}
{{end}}
`))
