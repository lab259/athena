package config

import (
	"os"
	"path"

	rscsrv "github.com/lab259/go-rscsrv"
)

// ProjectRoot return project root (from "PROJECT_ROOT" or current working directory).
func ProjectRoot() string {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		if dir, err := os.Getwd(); err == nil {
			return dir
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

var defaultConfigurationLoader rscsrv.ConfigurationLoader
var defaultConfigurationUnmarshaler rscsrv.ConfigurationUnmarshaler

func init() {
	defaultConfigurationUnmarshaler = &rscsrv.ConfigurationUnmarshalerYaml{}
}

// Load loads a file to a pointer on the current environment loaded from the
// `ENV` environment variable.
//
// `file` defines which file should be loaded from the environment.
// `dst` is a pointer.
func Load(file string, dst interface{}) error {
	defaultConfigurationLoader = rscsrv.NewFileConfigurationLoader(configurationFolder())

	config, err := defaultConfigurationLoader.Load(file)
	if err != nil {
		return err
	}
	return defaultConfigurationUnmarshaler.Unmarshal(config, dst)
}
