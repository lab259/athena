package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var FindServiceTemplate = template.New("find_service.go", `package {{.Table}}

import (
	"context"

	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"gopkg.in/src-d/go-kallax.v1"
	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/errors"
)

// FindInput holds input information for Find service
type FindInput struct {
	{{.Model}}ID kallax.ULID
}

// FindOutput holds the output information from Find service
type FindOutput struct {
	{{.Model}} *models.{{.Model}}
}

// Find returns a {{.Model}}
func Find(ctx context.Context, input *FindInput) (*FindOutput, error) {
	db, err := psqlrscsrv.DefaultPsqlService.DB()
	if err != nil {
		return nil, err
	}

	store := models.New{{.Model}}Store(db)
	obj, err := store.FindOne(models.New{{.Model}}Query().FindByID(input.{{.Model}}ID))
	if err != nil {
		return nil, errors.Wrap(err,errors.Code("find-failed"), errors.Module("{{.Table}}_service"))
	}

	return &FindOutput{
		{{.Model}}: obj,
	}, nil
}
`)

var FindServiceTestTemplate = template.New("find_test.go", `package {{.Table}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Table}}"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/lab259/athena/testing/psqltest"
	"github.com/felipemfp/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/src-d/go-kallax.v1"
	"github.com/lab259/errors"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("Find", func() {
			
			rscsrvtest.Setup(psqltest.NewPsqlTestService())

			It("should find", func() {
				ctx := context.Background()

				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				store := models.New{{.Model}}Store(db)

				existing := models.New{{.Model}}()
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = kallax.NewULID()
				Expect(store.Insert(existing)).To(Succeed())

				input := {{.Table}}.FindInput{}
				input.{{.Model}}ID = existing.ID

				output, err := {{.Table}}.Find(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.{{$.Model}}.ID).To(Equal(existing.ID))
				{{- range .Fields}}
				Expect(output.{{$.Model}}.{{formatFieldName .}}).To(Equal(existing.{{formatFieldName .}}))
				{{- end}}
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Table}}.FindInput{}
				input.{{.Model}}ID = kallax.NewULID()

				output, err := {{.Table}}.Find(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(kallax.ErrNotFound))
			})
		})
	})
})
`)
