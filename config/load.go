package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
	rscsrv "github.com/lab259/go-rscsrv"
)

var (
	projectRoot string
)

// ProjectRoot return project root (from "PROJECT_ROOT" or current working directory).
func ProjectRoot() string {
	if projectRoot == "" {
		projectRoot = os.Getenv("PROJECT_ROOT")
		if dir, err := os.Getwd(); err == nil && projectRoot == "" {
			projectRoot = dir

			for {
				if _, err := os.Stat(path.Join(projectRoot, "go.sum")); err != nil {
					if projectRoot == "/" {
						projectRoot = "."
						break
					} else {
						projectRoot = path.Dir(projectRoot)
					}
				} else {
					break
				}
			}
		}
	}
	return projectRoot
}

// Environment return current environment (from "ENV").
func Environment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

func configurationFolder() string {
	projectRoot := ProjectRoot()
	if projectRoot == "" {
		return path.Join("configs", Environment())
	}
	return path.Join(projectRoot, "configs", Environment())
}

var (
	defaultConfigurationUnmarshaler rscsrv.ConfigurationUnmarshaler = &rscsrv.ConfigurationUnmarshalerYaml{}
)

// Load loads a file to a pointer on the current environment loaded from the
// `ENV` environment variable.
//
// `file` defines which file should be loaded from the environment.
// `dst` is a pointer.
func Load(file string, dst interface{}) error {
	loader := rscsrv.NewFileConfigurationLoader(configurationFolder())

	name := strings.TrimSuffix(file, filepath.Ext(file))

	config, err := loader.Load(file)
	if err != nil {
		return err
	}

	if err := defaultConfigurationUnmarshaler.Unmarshal(config, dst); err != nil {
		return err
	}

	return envconfig.Process(strcase.ToScreamingSnake(name), dst)
}
