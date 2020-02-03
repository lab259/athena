package fasthttprscsrv

import (
	"context"
	"fmt"

	"github.com/lab259/athena/config"
	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/valyala/fasthttp"
)

// FasthttpConfiguration TODO
type FasthttpConfiguration struct {
	Addr     string `yaml:"addr"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type enhancedFasthttpServer struct {
	server        *fasthttp.Server
	config        FasthttpConfiguration
	serviceConfig ServiceConfiguration
}

func (srv *enhancedFasthttpServer) listenAndServe() error {
	if srv.config.CertFile != "" && srv.config.KeyFile != "" {
		return srv.server.ListenAndServeTLS(srv.config.Addr, srv.config.CertFile, srv.config.KeyFile)
	}

	return srv.server.ListenAndServe(srv.config.Addr)
}

func (srv *enhancedFasthttpServer) StartWithContext(ctx context.Context) (err error) {
	done := make(chan bool, 1)

	go func() {
		<-ctx.Done()
		err = srv.server.Shutdown()
		close(done)
	}()

	if err := srv.server.ListenAndServe(srv.config.Addr); err != nil {
		return err
	}

	<-done
	return
}

func (srv *enhancedFasthttpServer) Name() string {
	return srv.serviceConfig.Name
}

func (srv *enhancedFasthttpServer) ApplyConfiguration(config interface{}) error {
	switch t := config.(type) {
	case *FasthttpConfiguration:
		srv.config = *t
	case FasthttpConfiguration:
		srv.config = t
	default:
		return rscsrv.ErrWrongConfigurationInformed
	}

	if srv.config.Addr == "" {
		return fmt.Errorf("addr not provided")
	}

	return nil
}

func (srv *enhancedFasthttpServer) LoadConfiguration() (interface{}, error) {
	var target FasthttpConfiguration
	if err := config.Load(srv.serviceConfig.Prefix, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// ServiceConfiguration holds additional configuration for the enhanced fasthttp.Server.
type ServiceConfiguration struct {
	Name   string
	Prefix string
}

// EnhancedFasthttpServer TODO
type EnhancedFasthttpServer interface {
	rscsrv.Service
	rscsrv.Configurable
	rscsrv.StartableWithContext
}

// EnhanceFasthttpServer enhances a fasthttp.Server into rscsrv.StartableWithContext.
func EnhanceFasthttpServer(server *fasthttp.Server, config ServiceConfiguration) EnhancedFasthttpServer {
	if config.Prefix == "" {
		config.Prefix = "fasthttp"
	}

	if config.Name == "" {
		config.Name = "Fasthttp Server Service"
	}

	return &enhancedFasthttpServer{
		server:        server,
		serviceConfig: config,
	}
}
