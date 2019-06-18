package make

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/jinzhu/inflection"
	"github.com/lab259/athena/athena/util"
	"github.com/lab259/athena/config"
)

func Model(cmd *cli.Cmd) {
	var data modelTemplateData

	cmd.Spec = "[OPTIONS] MODEL FIELD..."

	cmd.StringArgPtr(&data.Model, "MODEL", "", "model name")
	cmd.StringsArgPtr(&data.Fields, "FIELD", []string{}, "field (eg.: Name:string)")
	cmd.BoolOptPtr(&data.WithRepository, "r repository", false, "with repository")
	cmd.BoolOptPtr(&data.WithCRUD, "crud", false, "with crud services")

	cmd.Action = func() {
		data.Model = strcase.ToCamel(data.Model)
		data.Collection = strcase.ToSnake(inflection.Plural(data.Model))
		if data.WithCRUD {
			data.WithRepository = true
		}

		projectRoot := config.ProjectRoot()
		data.Project = filepath.Base(projectRoot)

		modelsDir := path.Join(projectRoot, "models")

		err := os.MkdirAll(modelsDir, os.ModePerm)
		util.HandleError(err, "Unable to create models directory.")

		modelFile := fmt.Sprintf("%s.go", path.Join(modelsDir, strcase.ToSnake(data.Model)))

		content := bytes.NewBuffer([]byte{})
		err = mongoModelTemplate.Execute(content, &data)
		util.HandleError(err, "Unable to execute model template.")

		err = ioutil.WriteFile(modelFile, content.Bytes(), 0644)
		util.HandleError(err, "Unable to create model template.")

		if data.WithCRUD {
			serviceDir := path.Join(projectRoot, "services", data.Collection)
			err := os.MkdirAll(serviceDir, os.ModePerm)
			util.HandleError(err, "Unable to create services directory.")

			createContent := bytes.NewBuffer([]byte{})
			err = mongoModelCreateTemplate.Execute(createContent, &data)
			util.HandleError(err, "Unable to execute create service template.")
			err = ioutil.WriteFile(path.Join(serviceDir, "create.go"), createContent.Bytes(), 0644)
			util.HandleError(err, "Unable to create create service template.")

			updateContent := bytes.NewBuffer([]byte{})
			err = mongoModelUpdateTemplate.Execute(updateContent, &data)
			util.HandleError(err, "Unable to execute update service template.")
			err = ioutil.WriteFile(path.Join(serviceDir, "update.go"), updateContent.Bytes(), 0644)
			util.HandleError(err, "Unable to create update service template.")

			deleteContent := bytes.NewBuffer([]byte{})
			err = mongoModelDeleteTemplate.Execute(deleteContent, &data)
			util.HandleError(err, "Unable to execute delete service template.")
			err = ioutil.WriteFile(path.Join(serviceDir, "delete.go"), deleteContent.Bytes(), 0644)
			util.HandleError(err, "Unable to create delete service template.")

			listContent := bytes.NewBuffer([]byte{})
			err = mongoModelListTemplate.Execute(listContent, &data)
			util.HandleError(err, "Unable to execute list service template.")
			err = ioutil.WriteFile(path.Join(serviceDir, "list.go"), listContent.Bytes(), 0644)
			util.HandleError(err, "Unable to create list service template.")

			findContent := bytes.NewBuffer([]byte{})
			err = mongoModelFindTemplate.Execute(findContent, &data)
			util.HandleError(err, "Unable to execute find service template.")
			err = ioutil.WriteFile(path.Join(serviceDir, "find.go"), findContent.Bytes(), 0644)
			util.HandleError(err, "Unable to create find service template.")
		}

		fmt.Printf("%s was created.\n", modelFile)
	}
}

type modelTemplateData struct {
	Project        string
	Model          string
	Collection     string
	WithRepository bool
	WithCRUD       bool
	Fields         []string
}

func formatField(field string) (string, error) {
	values := strings.Split(field, ":")
	return strings.Join(values[:2], " "), nil
}

func formatFieldOptional(field string) (string, error) {
	values := strings.Split(field, ":")
	if len(values) < 2 {
		return "", fmt.Errorf("invalid format for field: %s", field)
	}
	return values[0] + " *" + values[1], nil
}

func formatFieldName(field string) string {
	values := strings.Split(field, ":")
	return values[0]
}

func formatFieldTag(field string) (string, error) {
	values := strings.Split(field, ":")
	if len(values) < 2 {
		return "", fmt.Errorf("invalid format for field: %s", field)
	}
	return strings.ReplaceAll(strcase.ToLowerCamel(values[0]), "ID", "Id"), nil
}

func formatValidation(field string) string {
	values := strings.Split(field, ":")
	if len(values) < 3 {
		return "-"
	}
	return values[2]
}

func hasValidation(field string) bool {
	return len(strings.Split(field, ":")) > 2
}

var fieldFunctions = template.FuncMap{
	"formatField":         formatField,
	"formatFieldTag":      formatFieldTag,
	"formatValidation":    formatValidation,
	"hasValidation":       hasValidation,
	"formatFieldName":     formatFieldName,
	"formatFieldOptional": formatFieldOptional,
}

var mongoModelTemplate = template.Must(template.New("make:model:mongo").Funcs(fieldFunctions).Parse(`package models

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
	{{range .Fields}}{{formatField .}}  ` + "`" + `json:"{{formatFieldTag .}}" bson:"{{formatFieldTag .}}"{{if hasValidation .}} validate:"{{formatValidation .}}"{{end}}` + "`" + `
	{{end}}
}

{{if .WithRepository}}
// New{{.Model}}Repository returns a Repository instance for {{.Model}} model
func New{{.Model}}Repository(ctx context.Context) *repository.Repository {
	return repository.NewRepository(repository.RepositoryConfig{
		Collection:  "{{.Collection}}",
  		QueryRunner: &mgo.DefaultMgoService,
	})
}
{{end}}
`))

var mongoModelCreateTemplate = template.Must(template.New("make:model:mongo:create").Funcs(fieldFunctions).Parse(`package {{.Collection}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
	"github.com/gofrs/uuid"
)

// CreateInput holds input information for Create service
type CreateInput struct {
	{{range .Fields}}{{formatField .}}  ` + "`" + `json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"{{formatValidation .}}"{{end}}` + "`" + `
	{{end}}
}

// CreateOutput holds the output information from Create service
type CreateOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Create creates a new {{.Model}}
func Create(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	repo := models.New{{.Model}}Repository(ctx)

	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("users_service"))
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("uuid-generation-failed"), errors.Module("users_service"))
	}

	obj := models.{{.Model}}{
		ID: uid,
		{{range .Fields}}{{formatFieldName .}}: input.{{formatFieldName .}},
		{{end}}
	}

	err = repo.Create(&obj)
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("repository-create-failed"), errors.Module("users_service"))
	}

	return &CreateOutput{
		{{.Model}}: &obj,
	}, nil
}
`))
var mongoModelUpdateTemplate = template.Must(template.New("make:model:mongo:update").Funcs(fieldFunctions).Parse(`package {{.Collection}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
	"github.com/lab259/mgohelpers"
	"github.com/gofrs/uuid"
)

// UpdateInput holds input information for Update service
type UpdateInput struct {
	{{.Model}}ID uuid.UUID
	{{range .Fields}}{{formatFieldOptional .}}  ` + "`" + `json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"omitempty,{{formatValidation .}}"{{end}}` + "`" + `
	{{end}}
}

// UpdateOutput holds the output information from Update service
type UpdateOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Update partial updates a {{.Model}} and returns it
func Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	var obj models.{{.Model}}
	repo := models.New{{.Model}}Repository(ctx)
	
	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("users_service"))
	}

	updateSet := make(mgohelpers.UpdateSet, {{len .Fields}})
	{{range .Fields}}updateSet.Add(&obj, &obj.{{formatFieldName .}}, input.{{formatFieldName .}})
	{{end}}

	err := repo.UpdateAndFind(input.{{.Model}}ID, &obj, updateSet)
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("repository-update-failed"), errors.Module("users_service"))
	}

	return &UpdateOutput{
		{{.Model}}: &obj,
	}, nil
}
`))
var mongoModelDeleteTemplate = template.Must(template.New("make:model:mongo:delete").Funcs(fieldFunctions).Parse(`package {{.Collection}}

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
	"github.com/lab259/repository"
	"github.com/gofrs/uuid"
)

// DeleteInput holds input information for Delete service
type DeleteInput struct {
	{{.Model}}ID uuid.UUID
}

// DeleteOutput holds the output information from Delete service
type DeleteOutput struct {
	Count int
}

// Delete deletes a {{.Model}}
func Delete(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	repo := models.New{{.Model}}Repository(ctx)

	err := repo.Delete(repository.ByID(input.{{.Model}}ID))
	if err != nil {
		return nil, errors.Wrap(err,errors.Code("repository-delete-failed"), errors.Module("users_service"))
	}

	return &DeleteOutput{
		Count: 1,
	}, nil
}
`))
var mongoModelListTemplate = template.Must(template.New("make:model:mongo:list").Funcs(fieldFunctions).Parse(`package {{.Collection}}

import (
	"context"

	"github.com/lab259/repository"
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/athena/pagination"
	"github.com/lab259/errors"
)

// ListInput holds input information for List service
type ListInput struct {
	CurrentPage int
	PageSize    int
}

// ListOutput holds the output information from List service
type ListOutput struct {
	Items       []*models.{{.Model}}
	Total       int
	CurrentPage int
	PageSize    int
}

// List returns a paginated list of {{.Model}}
func List(ctx context.Context, input *ListInput) (*ListOutput, error) {
	var objs []*models.{{.Model}}
	repo := models.New{{.Model}}Repository(ctx)

	pageSize, currentPage := pagination.Parse(input.PageSize, input.CurrentPage)

	total, err := repo.CountAndFindAll(&objs, repository.WithPage(currentPage-1, pageSize))
	if err != nil {
		return nil, errors.Wrap(err,errors.Code("repository-find-failed"), errors.Module("users_service"))
	}

	return {
		Items: objs,
		Total: total,
		CurrentPage: currentPage,
		PageSize: pageSize,
	}, nil
}
`))
var mongoModelFindTemplate = template.Must(template.New("make:model:mongo:find").Funcs(fieldFunctions).Parse(`package {{.Collection}}

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
	"github.com/lab259/repository"
	"github.com/gofrs/uuid"
)

// FindInput holds input information for Find service
type FindInput struct {
	{{.Model}}ID uuid.UUID
}

// FindOutput holds the output information from Find service
type FindOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Find returns a {{.Model}}
func Find(ctx context.Context, input *FindInput) (*FindOutput, error) {
	repo := models.New{{.Model}}Repository(ctx)
	var obj models.{{.Model}}

	err := repo.Find(&obj, repository.ByID(input.{{.Model}}ID))
	if err != nil {
		return nil, errors.Wrap(err,errors.Code("repository-find-failed"), errors.Module("users_service"))
	}

	return &FindOutput{
		{{.Model}}: &obj,
	}, nil
}
`))
