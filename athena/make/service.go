package make

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/lab259/athena/athena/make/templates"
	"github.com/lab259/athena/athena/util"
	"github.com/lab259/athena/athena/util/template"
	"github.com/lab259/athena/config"
)

func Service(cmd *cli.Cmd) {
	var data templates.ServiceTemplateData

	cmd.Spec = "PACKAGE SERVICE INPUT_FIELD..."

	cmd.StringArgPtr(&data.Package, "PACKAGE", "", "package name")
	cmd.StringArgPtr(&data.Service, "SERVICE", "", "service name")
	cmd.StringsArgPtr(&data.Fields, "INPUT_FIELD", []string{}, "input field (eg.: Name:string)")

	cmd.Action = func() {
		data.Service = strcase.ToCamel(data.Service)
		data.Package = strcase.ToSnake(data.Package)

		projectRoot := config.ProjectRoot()
		data.Project = filepath.Base(projectRoot)

		packageDir := path.Join(projectRoot, "services", data.Package)

		err := os.MkdirAll(packageDir, os.ModePerm)
		util.HandleError(err, "Unable to create package directory.")

		serviceFile := fmt.Sprintf("%s.go", path.Join(packageDir, strcase.ToSnake(data.Service)))
		serviceTestFile := fmt.Sprintf("%s_test.go", path.Join(packageDir, strcase.ToSnake(data.Service)))

		err = template.Write(templates.ServiceTemplate, &data, serviceFile)
		util.HandleError(err, "Unable to create service.")

		err = template.Write(templates.ServiceTestTemplate, &data, serviceTestFile)
		util.HandleError(err, "Unable to create service tests.")

		fmt.Println("The following files were created:")
		fmt.Printf("  %s\n", serviceFile)
		fmt.Printf("  %s\n", serviceTestFile)
	}
}
