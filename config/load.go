package config

import (
	"fmt"
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
	yamlConfigurationUnmarshaler rscsrv.ConfigurationUnmarshaler = &rscsrv.ConfigurationUnmarshalerYaml{}
)

// Load loads a file to a pointer on the current environment loaded from the
// `ENV` environment variable.
//
// `file` defines which file should be loaded from the environment.
// `dst` is a pointer.
func Load(fileOrPrefix string, dst interface{}) error {
	loader := rscsrv.NewFileConfigurationLoader(configurationFolder())

	ext := filepath.Ext(fileOrPrefix)
	name := strings.TrimSuffix(fileOrPrefix, ext)

	if config, err := loader.Load(fileOrPrefix); err == nil {
		var err error
		switch ext {
		case ".yaml", ".yml":
			err = yamlConfigurationUnmarshaler.Unmarshal(config, dst)
		default:
			err = fmt.Errorf("unexpected file extension: %s", ext)
		}
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return envconfig.Process(strcase.ToScreamingSnake(name), dst)
}
