package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

var (
	version = "dev"
)

func main() {
	app := cli.App("athena", "Wisely building web applications")

	app.Version("v version", version)

	app.Run(os.Args)
}
