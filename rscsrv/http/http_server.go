package httprscsrv

import (
	"context"
	"fmt"
	"time"

	"github.com/lab259/athena/config"
	rscsrv "github.com/lab259/go-rscsrv"
	"net/http"
)

// HttpConfiguration TODO
type HttpConfiguration struct {
	Addr     string `yaml:"addr"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type enhanceHttpServer struct {
	server        *http.Server
	config        HttpConfiguration
	serviceConfig ServiceConfiguration
}

func (srv *enhanceHttpServer) listenAndServe() error {
	srv.server.Addr = srv.config.Addr

	if srv.config.CertFile != "" && srv.config.KeyFile != "" {

		return srv.server.ListenAndServeTLS(srv.config.CertFile, srv.config.KeyFile)
	}

	return srv.server.ListenAndServe()
}

func (srv *enhanceHttpServer) StartWithContext(ctx context.Context) (err error) {
	done := make(chan bool, 1)

	go func() {
		<-ctx.Done()

		// TODO: make this ctx timeout configurable
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*20)
		defer cancelShutdown()

		errShutdown := srv.server.Shutdown(shutdownCtx)
		if errShutdown != nil {
			err = errShutdown
		}
		close(done)
	}()

	if errListen := srv.server.ListenAndServe(); errListen != nil {
		err = errListen
		return
	}

	<-done
	return
}

func (srv *enhanceHttpServer) Name() string {
	return srv.serviceConfig.Name
}

func (srv *enhanceHttpServer) ApplyConfiguration(config interface{}) error {
	switch t := config.(type) {
	case *HttpConfiguration:
		srv.config = *t
	case HttpConfiguration:
		srv.config = t
	default:
		return rscsrv.ErrWrongConfigurationInformed
	}

	if srv.server.Addr == "" && srv.config.Addr == "" {
		return fmt.Errorf("addr not provided")
	}

	return nil
}

func (srv *enhanceHttpServer) LoadConfiguration() (interface{}, error) {
	var target HttpConfiguration
	if err := config.Load(srv.serviceConfig.Prefix, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// ServiceConfiguration holds additional configuration for the enhanced http.Server.
type ServiceConfiguration struct {
	Name   string
	Prefix string
}

// HttpServer TODO
type HttpServer interface {
	rscsrv.Service
	rscsrv.Configurable
	rscsrv.StartableWithContext
}

// NewHttpServer enhances a http.Server into rscsrv.StartableWithContext.
func NewHttpServer(server *http.Server, config ServiceConfiguration) HttpServer {
	if config.Prefix == "" {
		config.Prefix = "Http"
	}

	if config.Name == "" {
		config.Name = "Http Server Service"
	}

	return &enhanceHttpServer{
		server:        server,
		serviceConfig: config,
	}
}
