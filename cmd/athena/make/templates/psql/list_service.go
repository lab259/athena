package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var ListServiceTemplate = template.New("list_service.go", `package {{.Table}}

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/athena/pagination"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
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
	db, err := psqlrscsrv.DefaultPsqlService.DB()
	if err != nil {
		return nil, errors.Wrap(err, errors.Code("db-available"), errors.Module("{{.Table}}_service"))
	}

	var output ListOutput

	store := models.New{{.Model}}Store(db)
	page := pagination.Parse(input.PageSize, input.CurrentPage)

	err = store.Transaction(func(s *models.{{.Model}}Store) error {
		q := models.New{{.Model}}Query().Limit(uint64(page.Limit)).Offset(uint64(page.Offset))
		objs, err := s.FindAll(q)
		if err != nil {
			return errors.Wrap(err, errors.Code("list-failed"), errors.Module("{{.Table}}_service"))
		}
		count, err := s.Count(q)
		if err != nil {
			return errors.Wrap(err, errors.Code("count-failed"), errors.Module("{{.Table}}_service"))
		}

		output.CurrentPage = page.CurrentPage
		output.PageSize = page.PageSize
		output.Items = objs
		output.Total = int(count)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &output, nil
}
`)

var ListServiceTestTemplate = template.New("list_test.go", `package {{.Table}}_test

import (
	"context"

	"github.com/felipemfp/faker"
	"github.com/lab259/athena/models"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/athena/services/accounts"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/lab259/athena/testing/psqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/src-d/go-kallax.v1"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("List", func() {
			
			rscsrvtest.Setup(psqltest.NewPsqlTestService())

			It("should list (empty)", func() {
				ctx := context.Background()
				input := {{.Table}}.ListInput{}
				
				output, err := {{.Table}}.List(ctx, &input)
				Expect(err).ToNot(HaveOccurred())
				
				Expect(output.Items).To(BeEmpty())
				Expect(output.Total).To(Equal(0))
				Expect(output.CurrentPage).To(Equal(1))
				Expect(output.PageSize).To(Equal(10))
			})

			It("should list (with 3, first page)", func() {
				ctx := context.Background()

				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				store := models.New{{.Model}}Store(db)

				existing1 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing1)).To(Succeed())
				existing1.ID = kallax.NewULID()
				Expect(store.Insert(existing1)).To(Succeed())

				existing2 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing2)).To(Succeed())
				existing2.ID = kallax.NewULID()
				Expect(store.Insert(existing2)).To(Succeed())

				existing3 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing3)).To(Succeed())
				existing3.ID = kallax.NewULID()
				Expect(store.Insert(existing3)).To(Succeed())

				input := {{.Table}}.ListInput{
					PageSize: 2,
				}

				output, err := {{.Table}}.List(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.CurrentPage).To(Equal(1))
				Expect(output.PageSize).To(Equal(2))
				Expect(output.Total).To(Equal(3))

				Expect(output.Items[0].ID).To(Equal(existing1.ID))
				{{- range .Fields}}
				Expect(output.Items[0].{{formatFieldName .}}).To(Equal(existing1.{{formatFieldName .}}))
				{{- end}}
				
				Expect(output.Items[1].ID).To(Equal(existing2.ID))
				{{- range .Fields}}
				Expect(output.Items[1].{{formatFieldName .}}).To(Equal(existing2.{{formatFieldName .}}))
				{{- end}}
			})

			It("should list (with 3, second page)", func() {
				ctx := context.Background()

				db, err := psqlrscsrv.DefaultPsqlService.DB()
				Expect(err).ToNot(HaveOccurred())

				store := models.New{{.Model}}Store(db)

				existing1 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing1)).To(Succeed())
				existing1.ID = kallax.NewULID()
				Expect(store.Insert(existing1)).To(Succeed())
				
				existing2 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing2)).To(Succeed())
				existing2.ID = kallax.NewULID()
				Expect(store.Insert(existing2)).To(Succeed())

				existing3 := models.New{{.Model}}()
				Expect(faker.FakeData(&existing3)).To(Succeed())
				existing3.ID = kallax.NewULID()
				Expect(store.Insert(existing3)).To(Succeed())

				input := {{.Table}}.ListInput{
					CurrentPage: 2,
					PageSize: 2,
				}

				output, err := {{.Table}}.List(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.CurrentPage).To(Equal(2))
				Expect(output.PageSize).To(Equal(2))
				Expect(output.Total).To(Equal(3))

				Expect(output.Items[0].ID).To(Equal(existing3.ID))
				{{- range .Fields}}
				Expect(output.Items[0].{{formatFieldName .}}).To(Equal(existing3.{{formatFieldName .}}))
				{{- end}}
			})
		})
	})
})
`)
