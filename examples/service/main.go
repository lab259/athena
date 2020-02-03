package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lab259/athena"
	fasthttprscsrv "github.com/lab259/athena/rscsrv/fasthttp"
	hermesrscsrv "github.com/lab259/athena/rscsrv/hermes"
	"github.com/lab259/hermes"
	"github.com/valyala/fasthttp"
)

type serviceA struct{}

func (srv *serviceA) Start() error {
	time.Sleep(2 * time.Second)
	return nil
}

func (srv *serviceA) Stop() error {
	time.Sleep(2 * time.Second)
	return nil
}

func (srv *serviceA) Name() string {
	return "Service A"
}

func main() {
	// 1. application using fasthttp
	var server fasthttp.Server
	server.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("version: v0.1.0")
	}
	serverService := fasthttprscsrv.EnhanceFasthttpServer(&server, fasthttprscsrv.ServiceConfiguration{
		Prefix: "BIND",
	})

	// 2. application using hermes
	router := hermes.DefaultRouter()
	router.Get("/version", func(_ hermes.Request, res hermes.Response) hermes.Result {
		return res.Data(map[string]string{"version": "v0.1.0"})
	})
	hermesService := hermesrscsrv.EnhanceHermesRouter(router, hermesrscsrv.ServiceConfiguration{})

	// 3. create an project cli
	cli, opt := athena.NewCLI("serviceapp", `Service Example`).
		Version("v0.1.0", "da39a3ee5e6b4b0d3255bfef95601890afd80709").
		Simple().
		Build()

	// 4. (optional) add some logic to before action
	cli.Before = func() {
		fmt.Printf("Version: %s (%s)\nEnvironment: %s\n\n", opt.Version, opt.Build, opt.Environment)
	}

	// 5. add applications as commands (could add other commands too)
	cli.Command("app", "REST API Server", athena.NewCommand(&serviceA{}, serverService).Build())
	cli.Command("hermes", "Hermes Server", athena.NewCommand(&serviceA{}, hermesService).Build())

	// 6. run
	cli.Run(os.Args)
}
