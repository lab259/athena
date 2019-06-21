package make

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	templates_psql "github.com/lab259/athena/cmd/athena/make/templates/psql"

	"github.com/lab259/athena/cmd/athena/util/template"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/jinzhu/inflection"
	"github.com/lab259/athena/cmd/athena/util"
	"github.com/lab259/athena/config"
)

func Model(cmd *cli.Cmd) {
	var data templates_psql.ModelTemplateData

	cmd.Spec = "[OPTIONS] MODEL FIELD..."

	cmd.StringArgPtr(&data.Model, "MODEL", "", "model name")
	cmd.StringsArgPtr(&data.Fields, "FIELD", []string{}, "field (eg.: Name:string)")
	cmd.BoolOptPtr(&data.WithCRUD, "crud", false, "with crud services")

	cmd.Action = func() {
		createdFiles := make([]string, 0, 12)
		defer dumpFiles(&createdFiles)

		data.Model = strcase.ToCamel(data.Model)
		data.Table = strcase.ToSnake(inflection.Plural(data.Model))

		projectRoot := config.ProjectRoot()
		data.Project = filepath.Base(projectRoot)

		modelsDir := path.Join(projectRoot, "models")

		err := os.MkdirAll(modelsDir, os.ModePerm)
		util.HandleError(err, "Unable to create models directory.")

		modelFile := fmt.Sprintf("%s.go", path.Join(modelsDir, strcase.ToSnake(data.Model)))
		err = template.Write(templates_psql.ModelTemplate, &data, modelFile)
		util.HandleError(err, "Unable to create model template.")
		createdFiles = append(createdFiles, modelFile)

		if data.WithCRUD {
			serviceDir := path.Join(projectRoot, "services", data.Table)
			err := os.MkdirAll(serviceDir, os.ModePerm)
			util.HandleError(err, "Unable to create services directory.")

			serviceTestsFile := path.Join(serviceDir, fmt.Sprintf("%s_test.go", data.Table))
			err = template.Write(templates_psql.ServiceTestsTemplate, &data, serviceTestsFile)
			util.HandleError(err, "Unable to create service tests.")
			createdFiles = append(createdFiles, serviceTestsFile)

			createServiceFile := path.Join(serviceDir, "create.go")
			err = template.Write(templates_psql.CreateServiceTemplate, &data, createServiceFile)
			util.HandleError(err, "Unable to create create service.")
			createdFiles = append(createdFiles, createServiceFile)
			createServiceTestFile := path.Join(serviceDir, "create_test.go")
			err = template.Write(templates_psql.CreateServiceTestTemplate, &data, createServiceTestFile)
			util.HandleError(err, "Unable to create create service tests.")
			createdFiles = append(createdFiles, createServiceTestFile)

			updateServiceFile := path.Join(serviceDir, "update.go")
			err = template.Write(templates_psql.UpdateServiceTemplate, &data, updateServiceFile)
			util.HandleError(err, "Unable to create update service.")
			createdFiles = append(createdFiles, updateServiceFile)
			updateServiceTestFile := path.Join(serviceDir, "update_test.go")
			err = template.Write(templates_psql.UpdateServiceTestTemplate, &data, updateServiceTestFile)
			util.HandleError(err, "Unable to create update service tests.")
			createdFiles = append(createdFiles, updateServiceTestFile)

			deleteServiceFile := path.Join(serviceDir, "delete.go")
			err = template.Write(templates_psql.DeleteServiceTemplate, &data, deleteServiceFile)
			util.HandleError(err, "Unable to create delete service.")
			createdFiles = append(createdFiles, deleteServiceFile)
			deleteServiceTestFile := path.Join(serviceDir, "delete_test.go")
			err = template.Write(templates_psql.DeleteServiceTestTemplate, &data, deleteServiceTestFile)
			util.HandleError(err, "Unable to create delete service tests.")
			createdFiles = append(createdFiles, deleteServiceTestFile)

			listServiceFile := path.Join(serviceDir, "list.go")
			err = template.Write(templates_psql.ListServiceTemplate, &data, listServiceFile)
			util.HandleError(err, "Unable to create list service.")
			createdFiles = append(createdFiles, listServiceFile)
			listServiceTestFile := path.Join(serviceDir, "list_test.go")
			err = template.Write(templates_psql.ListServiceTestTemplate, &data, listServiceTestFile)
			util.HandleError(err, "Unable to create list service tests.")
			createdFiles = append(createdFiles, listServiceTestFile)

			findServiceFile := path.Join(serviceDir, "find.go")
			err = template.Write(templates_psql.FindServiceTemplate, &data, findServiceFile)
			util.HandleError(err, "Unable to create find service.")
			createdFiles = append(createdFiles, findServiceFile)
			findServiceTestFile := path.Join(serviceDir, "find_test.go")
			err = template.Write(templates_psql.FindServiceTestTemplate, &data, findServiceTestFile)
			util.HandleError(err, "Unable to create find service tests.")
			createdFiles = append(createdFiles, findServiceTestFile)
		}
	}
}
