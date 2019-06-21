package templates_psql

import "github.com/lab259/athena/cmd/athena/util/template"

var UpdateServiceTemplate = template.New("update_service.go", `package {{.Table}}

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

var UpdateServiceTestTemplate = template.New("update_test.go", `package {{.Table}}_test

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
		Describe("Update", func() {
			
			BeforeEach(func() {
				rscsrvtest.Start(&mgorscsrv.DefaultMgoService)
				mgotest.ClearDefaultMgoService("")
			})

			It("should update", func() {
				ctx := context.Background()
				repo := models.New{{.Model}}Repository(ctx)

				existing := models.{{.Model}}{}
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = uuid.Must(uuid.NewV4())
				repo.Create(&existing)

				input := {{.Table}}.UpdateInput{}
				Expect(faker.FakeData(&input)).To(Succeed())
				input.{{.Model}}ID = existing.ID

				output, err := {{.Table}}.Update(ctx, &input)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.{{$.Model}}.ID).To(Equal(existing.ID))
				{{range .Fields}}Expect(output.{{$.Model}}.{{formatFieldName .}}).To(Equal(*input.{{formatFieldName .}}))
				{{end}}

				var obj models.{{.Model}}
				Expect(repo.FindByID(output.{{.Model}}.ID, &obj)).To(Succeed())

				{{range .Fields}}Expect(obj.{{formatFieldName .}}).To(Equal(*input.{{formatFieldName .}}))
				{{end}}
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Table}}.UpdateInput{}
				Expect(faker.FakeData(&input)).To(Succeed())
				input.{{.Model}}ID = uuid.Must(uuid.NewV4())

				output, err := {{.Table}}.Update(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(mgo.ErrNotFound))
			})
		})
	})
})
`)
