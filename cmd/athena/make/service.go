package make

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/lab259/athena/cmd/athena/make/templates"
	"github.com/lab259/athena/cmd/athena/util"
	"github.com/lab259/athena/cmd/athena/util/template"
	"github.com/lab259/athena/config"
)

func Service(cmd *cli.Cmd) {
	var data templates.ServiceTemplateData

	cmd.Spec = "PACKAGE SERVICE INPUT_FIELD..."

	cmd.StringArgPtr(&data.Package, "PACKAGE", "", "package name")
	cmd.StringArgPtr(&data.Service, "SERVICE", "", "service name")
	cmd.StringsArgPtr(&data.Fields, "INPUT_FIELD", []string{}, "input field (eg.: Name:string)")

	cmd.Action = func() {
		createdFiles := make([]string, 0, 3)
		defer dumpFiles(&createdFiles)

		data.Service = strcase.ToCamel(data.Service)
		data.Package = strcase.ToSnake(data.Package)

		projectRoot := config.ProjectRoot()
		data.Project = filepath.Base(projectRoot)

		packageDir := path.Join(projectRoot, "services", data.Package)

		err := os.MkdirAll(packageDir, os.ModePerm)
		util.HandleError(err, "Unable to create package directory.")

		serviceFile := fmt.Sprintf("%s.go", path.Join(packageDir, strcase.ToSnake(data.Service)))
		serviceTestFile := fmt.Sprintf("%s_test.go", path.Join(packageDir, strcase.ToSnake(data.Service)))
		packageTestFile := fmt.Sprintf("%s_test.go", path.Join(packageDir, data.Package))

		err = template.Write(templates.ServiceTemplate, &data, serviceFile)
		util.HandleError(err, "Unable to create service.")
		createdFiles = append(createdFiles, serviceFile)

		err = template.Write(templates.ServiceTestTemplate, &data, serviceTestFile)
		util.HandleError(err, "Unable to create service tests.")
		createdFiles = append(createdFiles, serviceTestFile)

		if _, err := os.Stat(packageTestFile); os.IsNotExist(err) {
			err = template.Write(templates.ServiceTestsTemplate, &data, packageTestFile)
			util.HandleError(err, "Unable to create package tests.")
			createdFiles = append(createdFiles, packageTestFile)
		}
	}
}

func dumpFiles(files *[]string) {
	fmt.Println("The following files were created:")

	dir, err := os.Getwd()
	for _, file := range *files {
		if err != nil {
			fmt.Printf("  %s\n", file)
		} else {
			fmt.Printf("  %s\n", strings.Replace(file, dir, ".", 1))
		}
	}
}
