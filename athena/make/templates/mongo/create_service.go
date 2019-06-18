package templates_mongo

import "github.com/lab259/athena/athena/util/template"

var CreateServiceTemplate = template.New("create_service.go", `package {{.Collection}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
	"github.com/gofrs/uuid"
)

// CreateInput holds input information for Create service
type CreateInput struct {
	{{range .Fields}}{{formatField .}}  `+"`"+`json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"{{formatValidation .}}"{{end}}`+"`"+`
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
`)
