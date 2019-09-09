package main

import (
	"fmt"
	"os"

	"github.com/lab259/athena"
)

func main() {
	cli, opt := athena.NewCLI("multiapp", `Multiapp Example`).
		Version("v0.1.0", "da39a3ee5e6b4b0d3255bfef95601890afd80709").
		Simple().
		Build()

	cli.Before = func() {
		fmt.Printf("Version: %s (%s)\nEnvironment: %s\n\n", opt.Version, opt.Build, opt.Environment)
	}

	cli.Command("app", "REST API Server", athena.NewCommand(func(opt *athena.CommandOptions) {
		fmt.Println("app server listening to " + opt.BindAddress)
	}).Build())

	cli.Command("graphql", "GraphQL Server", athena.NewCommand(func(opt *athena.CommandOptions) {
		fmt.Println("graphql server listening to " + opt.BindAddress)
	}).Build())

	cli.Run(os.Args)
}
