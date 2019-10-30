package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var DeleteServiceTemplate = template.New("delete_service.go", `package {{.Table}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/{{.Project}}/models"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/errors"
)

// DeleteInput holds input information for Delete service
type DeleteInput struct {
	{{.Model}} *models.{{.Model}} `+"`"+`validate:"required,structonly"`+"`"+`
}

// DeleteOutput holds the output information from Delete service
type DeleteOutput struct {
	Count int
}

// Delete deletes a {{.Model}}
func Delete(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("{{.Table}}_service"))
	}

	db, err := psqlrscsrv.DefaultPsqlService.DB()
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("db-not-available"), errors.Module("{{.Table}}_service"))
	}

	store := models.New{{.Model}}Store(db)
	err = store.Delete(input.{{.Model}})
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("delete-failed"), errors.Module("{{.Table}}_service"))
	}

	return &DeleteOutput{
		Count: 1,
	}, nil
}
`)

var DeleteServiceTestTemplate = template.New("delete_test.go", `package {{.Table}}_test

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
		Describe("Delete", func() {

			rscsrvtest.Setup(psqltest.NewPsqlTestService())

			It("should delete", func() {
				ctx := context.Background()

				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				store := models.New{{.Model}}Store(db)

				existing := models.New{{.Model}}()
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = kallax.NewULID()
				Expect(store.Insert(existing)).To(Succeed())

				input := {{.Table}}.DeleteInput{}
				input.{{.Model}} = existing

				output, err := {{.Table}}.Delete(ctx, &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Count).To(Equal(1))

				_, err = store.FindOne(models.New{{.Model}}Query().FindByID(existing.ID))
				Expect(err).To(HaveOccurred())
				Expect(errors.Reason(err)).To(Equal(kallax.ErrNotFound))
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Table}}.DeleteInput{}
				input.{{.Model}} = models.New{{.Model}}()
				Expect(faker.FakeData(&input.{{.Model}})).To(Succeed())

				output, err := {{.Table}}.Delete(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(kallax.ErrEmptyID))
			})
		})
	})
})
`)
