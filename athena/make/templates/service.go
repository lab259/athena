package templates

import "github.com/lab259/athena/athena/util/template"

type ServiceTemplateData struct {
	Project string
	Service string
	Package string
	Fields  []string
}

var ServiceTemplate = template.New("service.go", `package {{.Package}}

import (
	"context"

	"github.com/lab259/athena/validator"
	"github.com/lab259/errors"
)

// {{.Service}}Input holds input information for {{.Service}} service
type {{.Service}}Input struct {
	{{range .Fields}}{{formatField .}}  `+"`"+`json:"{{formatFieldTag .}}" {{if hasValidation .}}validate:"{{formatValidation .}}"{{end}}`+"`"+`
	{{end}}
}

// {{.Service}}Output holds the output information from {{.Service}} service
type {{.Service}}Output struct {
	// TODO
}

// {{.Service}} TODO
func {{.Service}}(ctx context.Context, input *{{.Service}}Input) (*{{.Service}}Output, error) {
	err := validator.Validate(input)
	if err != nil {
		return nil, errors.Wrap(err, errors.Validation(), errors.Module("{{.Package}}_service"))
	}

	panic(errors.New("not implemented"))
}
`)

var ServiceTestTemplate = template.New("service_test.go", `package {{.Package}}_test

import (
	"context"

	"github.com/lab259/{{.Project}}/services/{{.Package}}"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/lab259/errors"
)

var _ = Describe("Services", func() {
	Describe("{{toCamel .Package}}", func() {
		Describe("{{.Service}}", func() {
			
			PIt("TODO", func() {
				ctx := context.Background()

				input := {{.Package}}.{{.Service}}Input{}

				output, err := {{.Package}}.{{.Service}}(ctx, &input)
				
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
			})
		})
	})
})
`)
