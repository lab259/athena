package templates_mongo

import "github.com/lab259/athena/athena/util/template"

var UpdateServiceTemplate = template.New("update_service.go", `package {{.Collection}}

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
	{{range .Fields}}{{formatFieldOptional .}}  `+"`"+`json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"omitempty,{{formatValidation .}}"{{end}}`+"`"+`
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

	err = repo.UpdateAndFind(input.{{.Model}}ID, &obj, updateSet)
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("repository-update-failed"), errors.Module("users_service"))
	}

	return &UpdateOutput{
		{{.Model}}: &obj,
	}, nil
}
`)

var UpdateServiceTestTemplate = template.New("update_test.go", `package {{.Collection}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/services/{{.Collection}}"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Collection}}", func() {
		Describe("Update", func() {
			
			PIt("TODO", func() {
				ctx := context.Background()

				input := {{.Collection}}.UpdateInput{}

				output, err := {{.Collection}}.Update(ctx, &input)
				
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
			})
		})
	})
})
`)
