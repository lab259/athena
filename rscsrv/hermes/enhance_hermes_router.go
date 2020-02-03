package hermesrscsrv

import (
	fasthttprscsrv "github.com/lab259/athena/rscsrv/fasthttp"
	"github.com/lab259/go-rscsrv"
	"github.com/lab259/hermes"
	"github.com/valyala/fasthttp"
)

// ServiceConfiguration holds additional configuration for the enhanced hermes.Router.
type ServiceConfiguration struct {
	Name   string
	Prefix string
}

// EnhancedHermesRouter TODO
type EnhancedHermesRouter interface {
	rscsrv.Service
	rscsrv.Configurable
	rscsrv.StartableWithContext
}

// EnhanceHermesRouter TODO
func EnhanceHermesRouter(router hermes.Router, config ServiceConfiguration) EnhancedHermesRouter {
	server := fasthttp.Server{
		Handler: router.Handler(),
	}

	if config.Name == "" {
		config.Name = "Hermes Server Service"
	}

	if config.Prefix == "" {
		config.Prefix = "hermes"
	}

	return fasthttprscsrv.EnhanceFasthttpServer(&server, fasthttprscsrv.ServiceConfiguration{
		Name:   config.Name,
		Prefix: config.Prefix,
	})
}
