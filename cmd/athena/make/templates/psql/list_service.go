package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var ListServiceTemplate = template.New("list_service.go", `package {{.Table}}

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
		return nil, errors.Wrap(err,errors.Code("repository-list-failed"), errors.Module("users_service"))
	}

	return &ListOutput{
		Items: objs,
		Total: total,
		CurrentPage: currentPage,
		PageSize: pageSize,
	}, nil
}
`)

var ListServiceTestTemplate = template.New("list_test.go", `package {{.Table}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Table}}"
	mgorscsrv "github.com/lab259/athena/rscsrv/mgo"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/lab259/athena/testing/mgotest"
	"github.com/gofrs/uuid"
	"github.com/felipemfp/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("List", func() {
			
			BeforeEach(func() {
				rscsrvtest.Start(&mgorscsrv.DefaultMgoService)
				mgotest.ClearDefaultMgoService("")
			})

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
				repo := models.New{{.Model}}Repository(ctx)

				existing1 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing1)).To(Succeed())
				existing1.ID = uuid.FromStringOrNil("397336c5-b8cd-4581-97cd-ba03568c5191")
				Expect(repo.Create(&existing1)).To(Succeed())

				existing2 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing2)).To(Succeed())
				existing2.ID = uuid.FromStringOrNil("406ff0d8-3af7-43b6-a595-17969d6def71")
				Expect(repo.Create(&existing2)).To(Succeed())

				existing3 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing3)).To(Succeed())
				existing3.ID = uuid.FromStringOrNil("82818f4d-a4be-4ee9-99f3-f7f4b9cdd910")
				Expect(repo.Create(&existing3)).To(Succeed())

				input := {{.Table}}.ListInput{
					PageSize: 2,
				}

				output, err := {{.Table}}.List(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.CurrentPage).To(Equal(1))
				Expect(output.PageSize).To(Equal(2))
				Expect(output.Total).To(Equal(3))

				Expect(output.Items[0].ID).To(Equal(existing1.ID))
				{{range .Fields}}Expect(output.Items[0].{{formatFieldName .}}).To(Equal(existing1.{{formatFieldName .}}))
				{{end}}
				Expect(output.Items[1].ID).To(Equal(existing2.ID))
				{{range .Fields}}Expect(output.Items[1].{{formatFieldName .}}).To(Equal(existing2.{{formatFieldName .}}))
				{{end}}
			})

			It("should list (with 3, second page)", func() {
				ctx := context.Background()
				repo := models.New{{.Model}}Repository(ctx)

				existing1 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing1)).To(Succeed())
				existing1.ID = uuid.FromStringOrNil("397336c5-b8cd-4581-97cd-ba03568c5191")
				Expect(repo.Create(&existing1)).To(Succeed())

				existing2 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing2)).To(Succeed())
				existing2.ID = uuid.FromStringOrNil("406ff0d8-3af7-43b6-a595-17969d6def71")
				Expect(repo.Create(&existing2)).To(Succeed())

				existing3 := models.{{.Model}}{}
				Expect(faker.FakeData(&existing3)).To(Succeed())
				existing3.ID = uuid.FromStringOrNil("82818f4d-a4be-4ee9-99f3-f7f4b9cdd910")
				Expect(repo.Create(&existing3)).To(Succeed())

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
				{{range .Fields}}Expect(output.Items[0].{{formatFieldName .}}).To(Equal(existing3.{{formatFieldName .}}))
				{{end}}
			})
		})
	})
})
`)
