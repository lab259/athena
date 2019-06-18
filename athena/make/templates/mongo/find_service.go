package templates_mongo

import "github.com/lab259/athena/athena/util/template"

var FindServiceTemplate = template.New("find_service.go", `package {{.Collection}}

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
`)