package make

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/lab259/athena/athena/util"
)

func Service(cmd *cli.Cmd) {
	var data serviceTemplateData

	cmd.StringArgPtr(&data.Package, "PACKAGE", "", "package name")
	cmd.StringArgPtr(&data.Service, "SERVICE", "", "service name")

	cmd.Action = func() {
		data.Service = strcase.ToCamel(data.Service)
		data.Package = strcase.ToSnake(data.Package)

		dir, err := os.Getwd()
		util.HandleError(err, "Unable to get current directory.")

		servicesDir := path.Join(dir, "services")
		packageDir := path.Join(servicesDir, data.Package)

		err = os.MkdirAll(packageDir, os.ModePerm)
		util.HandleError(err, "Unable to create package directory.")

		serviceFile := fmt.Sprintf("%s.go", path.Join(packageDir, strcase.ToSnake(data.Service)))

		content := bytes.NewBuffer([]byte{})
		err = serviceTemplate.Execute(content, &data)
		util.HandleError(err, "Unable to execute service template.")

		err = ioutil.WriteFile(serviceFile, content.Bytes(), 0644)
		util.HandleError(err, "Unable to create service template.")

		fmt.Printf("%s was created.\n", serviceFile)
	}
}

type serviceTemplateData struct {
	Service string
	Package string
}

var serviceTemplate = template.Must(template.New("make:service").Parse(`package {{.Package}}

import (
	"context"
	"github.com/lab259/errors"
)

// {{.Service}}Input holds input information for {{.Service}} service
type {{.Service}}Input struct {
	// TODO
}

// {{.Service}}Output holds the output information from {{.Service}} service
type {{.Service}}Output struct {
	// TODO
}

// {{.Service}} TODO
func {{.Service}}(ctx context.Context, input *{{.Service}}Input) (*{{.Service}}Output, error) {
	return nil, errors.New("not implemented")
}
`))
