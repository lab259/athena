package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var FindServiceTemplate = template.New("find_service.go", `package {{.Table}}

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

var FindServiceTestTemplate = template.New("find_test.go", `package {{.Table}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Table}}"
	mgorscsrv "github.com/lab259/athena/rscsrv/mgo"
	"github.com/lab259/athena/testing/rscsrvtest"
	"github.com/lab259/athena/testing/mgotest"
	"github.com/gofrs/uuid"
	"github.com/globalsign/mgo"
	"github.com/lab259/errors"
	"github.com/felipemfp/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Table}}", func() {
		Describe("Find", func() {
			
			BeforeEach(func() {
				rscsrvtest.Start(&mgorscsrv.DefaultMgoService)
				mgotest.ClearDefaultMgoService("")
			})

			It("should find", func() {
				ctx := context.Background()
				repo := models.New{{.Model}}Repository(ctx)

				existing := models.{{.Model}}{}
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = uuid.Must(uuid.NewV4())
				repo.Create(&existing)

				input := {{.Table}}.FindInput{}
				input.{{.Model}}ID = existing.ID

				output, err := {{.Table}}.Find(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.{{$.Model}}.ID).To(Equal(existing.ID))
				{{range .Fields}}Expect(output.{{$.Model}}.{{formatFieldName .}}).To(Equal(existing.{{formatFieldName .}}))
				{{end}}
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Table}}.FindInput{}
				input.{{.Model}}ID = uuid.Must(uuid.NewV4())

				output, err := {{.Table}}.Find(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(mgo.ErrNotFound))
			})
		})
	})
})
`)
