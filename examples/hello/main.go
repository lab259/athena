package main

import (
	"fmt"
	"os"

	"github.com/lab259/athena"
)

func main() {
	cli, _ := athena.NewCLI("hello-world", `Prints "Hello, World!"`).
		Version("v0.1.0", "da39a3ee5e6b4b0d3255bfef95601890afd80709").
		Environment("development").
		Build()

	cli.Action = func() {
		fmt.Println("Hello, World!")
	}

	cli.Run(os.Args)
}
