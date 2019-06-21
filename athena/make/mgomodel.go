package make

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	templates_mongo "github.com/lab259/athena/athena/make/templates/mongo"

	"github.com/lab259/athena/athena/util/template"

	"github.com/iancoleman/strcase"
	cli "github.com/jawher/mow.cli"
	"github.com/jinzhu/inflection"
	"github.com/lab259/athena/athena/util"
	"github.com/lab259/athena/config"
)

func MgoModel(cmd *cli.Cmd) {
	var data templates_mongo.ModelTemplateData

	cmd.Spec = "[OPTIONS] MODEL FIELD..."

	cmd.StringArgPtr(&data.Model, "MODEL", "", "model name")
	cmd.StringsArgPtr(&data.Fields, "FIELD", []string{}, "field (eg.: Name:string)")
	cmd.BoolOptPtr(&data.WithRepository, "r repository", false, "with repository")
	cmd.BoolOptPtr(&data.WithCRUD, "crud", false, "with crud services")

	cmd.Action = func() {
		createdFiles := make([]string, 0, 12)
		defer dumpFiles(&createdFiles)

		data.Model = strcase.ToCamel(data.Model)
		data.Collection = strcase.ToSnake(inflection.Plural(data.Model))
		if data.WithCRUD {
			data.WithRepository = true
		}

		projectRoot := config.ProjectRoot()
		data.Project = filepath.Base(projectRoot)

		modelsDir := path.Join(projectRoot, "models")

		err := os.MkdirAll(modelsDir, os.ModePerm)
		util.HandleError(err, "Unable to create models directory.")

		modelFile := fmt.Sprintf("%s.go", path.Join(modelsDir, strcase.ToSnake(data.Model)))
		err = template.Write(templates_mongo.ModelTemplate, &data, modelFile)
		util.HandleError(err, "Unable to create model template.")
		createdFiles = append(createdFiles, modelFile)

		if data.WithCRUD {
			serviceDir := path.Join(projectRoot, "services", data.Collection)
			err := os.MkdirAll(serviceDir, os.ModePerm)
			util.HandleError(err, "Unable to create services directory.")

			serviceTestsFile := path.Join(serviceDir, fmt.Sprintf("%s_test.go", data.Collection))
			err = template.Write(templates_mongo.ServiceTestsTemplate, &data, serviceTestsFile)
			util.HandleError(err, "Unable to create service tests.")
			createdFiles = append(createdFiles, serviceTestsFile)

			createServiceFile := path.Join(serviceDir, "create.go")
			err = template.Write(templates_mongo.CreateServiceTemplate, &data, createServiceFile)
			util.HandleError(err, "Unable to create create service.")
			createdFiles = append(createdFiles, createServiceFile)
			createServiceTestFile := path.Join(serviceDir, "create_test.go")
			err = template.Write(templates_mongo.CreateServiceTestTemplate, &data, createServiceTestFile)
			util.HandleError(err, "Unable to create create service tests.")
			createdFiles = append(createdFiles, createServiceTestFile)

			updateServiceFile := path.Join(serviceDir, "update.go")
			err = template.Write(templates_mongo.UpdateServiceTemplate, &data, updateServiceFile)
			util.HandleError(err, "Unable to create update service.")
			createdFiles = append(createdFiles, updateServiceFile)
			updateServiceTestFile := path.Join(serviceDir, "update_test.go")
			err = template.Write(templates_mongo.UpdateServiceTestTemplate, &data, updateServiceTestFile)
			util.HandleError(err, "Unable to create update service tests.")
			createdFiles = append(createdFiles, updateServiceTestFile)

			deleteServiceFile := path.Join(serviceDir, "delete.go")
			err = template.Write(templates_mongo.DeleteServiceTemplate, &data, deleteServiceFile)
			util.HandleError(err, "Unable to create delete service.")
			createdFiles = append(createdFiles, deleteServiceFile)
			deleteServiceTestFile := path.Join(serviceDir, "delete_test.go")
			err = template.Write(templates_mongo.DeleteServiceTestTemplate, &data, deleteServiceTestFile)
			util.HandleError(err, "Unable to create delete service tests.")
			createdFiles = append(createdFiles, deleteServiceTestFile)

			listServiceFile := path.Join(serviceDir, "list.go")
			err = template.Write(templates_mongo.ListServiceTemplate, &data, listServiceFile)
			util.HandleError(err, "Unable to create list service.")
			createdFiles = append(createdFiles, listServiceFile)
			listServiceTestFile := path.Join(serviceDir, "list_test.go")
			err = template.Write(templates_mongo.ListServiceTestTemplate, &data, listServiceTestFile)
			util.HandleError(err, "Unable to create list service tests.")
			createdFiles = append(createdFiles, listServiceTestFile)

			findServiceFile := path.Join(serviceDir, "find.go")
			err = template.Write(templates_mongo.FindServiceTemplate, &data, findServiceFile)
			util.HandleError(err, "Unable to create find service.")
			createdFiles = append(createdFiles, findServiceFile)
			findServiceTestFile := path.Join(serviceDir, "find_test.go")
			err = template.Write(templates_mongo.FindServiceTestTemplate, &data, findServiceTestFile)
			util.HandleError(err, "Unable to create find service tests.")
			createdFiles = append(createdFiles, findServiceTestFile)
		}
	}
}
