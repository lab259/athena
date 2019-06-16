package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/lab259/athena/athena/make"
)

var (
	version = "dev"
)

func main() {
	app := cli.App("athena", "Wisely building web applications")

	app.Version("v version", version)

	app.Command("make:service", "Generate a service file", make.Service)

	app.Run(os.Args)
}
