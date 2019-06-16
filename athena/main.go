package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/lab259/athena/athena/make"
	"github.com/lab259/athena/athena/setup"
)

var (
	version = "dev"
)

func main() {
	app := cli.App("athena", "Wisely building web applications")

	app.Version("v version", version)

	app.Command("make:service", "Generate a service file", make.Service)
	app.Command("make:model", "Generate a model file", make.Model)

	app.Command("setup:sra", "Setup sra service", setup.Sra)

	app.Run(os.Args)
}
