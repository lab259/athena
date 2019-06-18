package templates_mongo

import "github.com/lab259/athena/athena/util/template"

var DeleteServiceTemplate = template.New("delete_service.go", `package {{.Collection}}

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
`)

var DeleteServiceTestTemplate = template.New("delete_test.go", `package {{.Collection}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/services/{{.Collection}}"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Collection}}", func() {
		Describe("Delete", func() {
			
			PIt("TODO", func() {
				ctx := context.Background()

				input := {{.Collection}}.DeleteInput{}

				output, err := {{.Collection}}.Delete(ctx, &input)
				
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
			})
		})
	})
})
`)
