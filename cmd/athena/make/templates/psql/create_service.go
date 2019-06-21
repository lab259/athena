package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var CreateServiceTemplate = template.New("create_service.go", `package {{.Table}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/errors"
)

// CreateInput holds input information for Create service
type CreateInput struct {
	{{- range .Fields}}
	{{formatField .}}  `+"`"+`json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"{{formatValidation .}}"{{end}}`+"`"+`
	{{- end}}
}

// CreateOutput holds the output information from Create service
type CreateOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Create creates a new {{.Model}}
func Create(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	db, err := psqlrscsrv.DefaultPsqlService.DB()
	if err != nil {
		return nil
	}

	store := models.New{{.Model}}Store(db)
	
	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("{{.Table}}_service"))
	}

	obj := models.New{{.Model}}()
	{{- range .Fields}}
	obj.{{formatFieldName .}} = input.{{formatFieldName .}}
	{{- end}}

	err = store.Insert(obj)
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("create-failed"), errors.Module("{{.Table}}_service"))
	}

	return &CreateOutput{
		{{.Model}}: &obj,
	}, nil
}
`)

var CreateServiceTestTemplate = template.New("create_test.go", `package {{.Table}}_test

import (
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Table}}"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/felipemfp/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("Create", func() {

			BeforeEach(func() {
				rscsrvtest.Start(&psqlrscsrv.DefaultPsqlService)
			})
			
			It("should create", func() {
				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				input := {{.Table}}.CreateInput{}
				Expect(faker.FakeData(&input)).To(Succeed())

				output, err := {{.Table}}.Create(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.{{$.Model}}.ID.IsEmpty()).To(BeFalse())
				{{- range .Fields}}
				Expect(output.{{$.Model}}.{{formatFieldName .}}).To(Equal(input.{{formatFieldName .}}))
				{{- end}}

				store := models.New{{.Model}}Store(db)
				obj, err := store.FindOne(models.New{{.Model}}Query().FindByID(output.{{.Model}}.ID))
				Expect(err).ToNot(HaveOccurred())

				{{- range .Fields}}
				Expect(obj.{{formatFieldName .}}).To(Equal(input.{{formatFieldName .}}))
				{{- end}}
			})
		})
	})
})
`)
