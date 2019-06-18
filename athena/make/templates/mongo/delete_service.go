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

	"github.com/lab259/{{.Project}}/models"
	"github.com/lab259/{{.Project}}/services/{{.Collection}}"
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
	Describe("{{toCamel .Collection}}", func() {
		Describe("Delete", func() {
			
			BeforeEach(func() {
				rscsrvtest.Start(&mgorscsrv.DefaultMgoService)
				mgotest.ClearDefaultMgoService("")
			})

			It("should delete", func() {
				ctx := context.Background()
				repo := models.New{{.Model}}Repository(ctx)

				existing := models.{{.Model}}{}
				Expect(faker.FakeData(&existing)).To(Succeed())
				existing.ID = uuid.Must(uuid.NewV4())
				repo.Create(&existing)

				input := {{.Collection}}.DeleteInput{}
				input.{{.Model}}ID = existing.ID

				output, err := {{.Collection}}.Delete(ctx, &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Count).To(Equal(1))

				var obj models.{{.Model}}
				err = repo.FindByID(existing.ID, &obj)
				Expect(err).To(HaveOccurred())
				Expect(errors.Reason(err)).To(Equal(mgo.ErrNotFound))
			})

			It("should fail with not found", func() {
				ctx := context.Background()

				input := {{.Collection}}.DeleteInput{}
				input.{{.Model}}ID = uuid.Must(uuid.NewV4())

				output, err := {{.Collection}}.Delete(ctx, &input)
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(errors.Reason(err)).To(Equal(mgo.ErrNotFound))
			})
		})
	})
})
`)
