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

func Model(cmd *cli.Cmd) {
	var data templates_mongo.ModelTemplateData

	cmd.Spec = "[OPTIONS] MODEL FIELD..."

	cmd.StringArgPtr(&data.Model, "MODEL", "", "model name")
	cmd.StringsArgPtr(&data.Fields, "FIELD", []string{}, "field (eg.: Name:string)")
	cmd.BoolOptPtr(&data.WithRepository, "r repository", false, "with repository")
	cmd.BoolOptPtr(&data.WithCRUD, "crud", false, "with crud services")

	cmd.Action = func() {
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

		if data.WithCRUD {
			serviceDir := path.Join(projectRoot, "services", data.Collection)
			err := os.MkdirAll(serviceDir, os.ModePerm)
			util.HandleError(err, "Unable to create services directory.")

			err = template.Write(templates_mongo.CreateServiceTemplate, &data, path.Join(serviceDir, "create.go"))
			util.HandleError(err, "Unable to create create service template.")

			err = template.Write(templates_mongo.UpdateServiceTemplate, &data, path.Join(serviceDir, "update.go"))
			util.HandleError(err, "Unable to create update service template.")

			err = template.Write(templates_mongo.DeleteServiceTemplate, &data, path.Join(serviceDir, "delete.go"))
			util.HandleError(err, "Unable to create delete service template.")

			err = template.Write(templates_mongo.ListServiceTemplate, &data, path.Join(serviceDir, "list.go"))
			util.HandleError(err, "Unable to create list service template.")

			err = template.Write(templates_mongo.FindServiceTemplate, &data, path.Join(serviceDir, "find.go"))
			util.HandleError(err, "Unable to create find service template.")
		}

		fmt.Printf("%s was created.\n", modelFile)
	}
}
