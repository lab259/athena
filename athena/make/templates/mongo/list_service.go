package templates_mongo

import "github.com/lab259/athena/athena/util/template"

var ListServiceTemplate = template.New("list_service.go", `package {{.Collection}}

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
`)
