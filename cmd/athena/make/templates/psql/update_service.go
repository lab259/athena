package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var UpdateServiceTemplate = template.New("update_service.go", `package {{.Table}}

import (
	"context"

	"gopkg.in/src-d/go-kallax.v1"
	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/errors"
)

// UpdateInput holds input information for Update service
type UpdateInput struct {
	{{.Model}} *models.{{.Model}}
	{{- range .Fields}}
	{{formatFieldOptional .}}  `+"`"+`json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"omitempty,{{formatValidation .}}"{{end}}`+"`"+`
	{{- end}}
}

// UpdateOutput holds the output information from Update service
type UpdateOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Update partial updates a {{.Model}} and returns it
func Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("{{.Table}}_service"))
	}

	db, err := psqlrscsrv.DefaultPsqlService.DB()
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("db-available"), errors.Module("{{.Table}}_service"))
	}

	store := models.New{{.Model}}Store(db)
	obj  := input.{{.Model}}
	cols := make([]kallax.SchemaField, 0, {{len .Fields}})

	{{range .Fields}}
	if input.{{formatFieldName .}} != nil {
		obj.{{formatFieldName .}} = *input.{{formatFieldName .}}
		cols = append(cols, models.Schema.{{$.Model}}.{{formatFieldName .}})
	}
	{{end}}

	_, err = store.Update(obj, cols...)
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("update-failed"), errors.Module("{{.Table}}_service"))
	}

	return &UpdateOutput{
		{{.Model}}: obj,
	}, nil
}
`)

var UpdateServiceTestTemplate = template.New("update_test.go", `package {{.Table}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Table}}"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/felipemfp/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/src-d/go-kallax.v1"
	"github.com/lab259/errors"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("Update", func() {

			BeforeEach(func() {
				rscsrvtest.Start(&psqlrscsrv.DefaultPsqlService)
			})

			AfterEach(func() {
				Expect(psqlrscsrv.DefaultPsqlService.Stop()).To(Succeed())
			})

			It("should update", func() {
				ctx := context.Background()

				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				store := models.New{{.Model}}Store(db)

				existing := models.New{{.Model}}()
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = kallax.NewULID()
				Expect(store.Insert(existing)).To(Succeed())

				input := {{.Table}}.UpdateInput{}
				Expect(faker.FakeData(&input)).To(Succeed())
				input.{{.Model}} = existing

				output, err := {{.Table}}.Update(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.{{$.Model}}.ID).To(Equal(existing.ID))
				{{- range .Fields}}
				Expect(output.{{$.Model}}.{{formatFieldName .}}).To(Equal(*input.{{formatFieldName .}}))
				{{- end}}

				obj, err := store.FindOne(models.New{{.Model}}Query().FindByID(output.{{.Model}}.ID))
				Expect(err).ToNot(HaveOccurred())

				{{- range .Fields}}
				Expect(obj.{{formatFieldName .}}).To(Equal(*input.{{formatFieldName .}}))
				{{- end}}
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Table}}.UpdateInput{}
				Expect(faker.FakeData(&input)).To(Succeed())
				input.{{.Model}} = models.New{{.Model}}()
				input.{{.Model}}.ID = kallax.NewULID()

				output, err := {{.Table}}.Update(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(kallax.ErrNoRowUpdate))
			})
		})
	})
})
`)
