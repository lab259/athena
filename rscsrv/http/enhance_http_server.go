package httprscsrv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lab259/athena/config"
	rscsrv "github.com/lab259/go-rscsrv"
)

type HTTPConfiguration struct {
	Addr     string `yaml:"addr"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type enhancedHTTPServer struct {
	server        *http.Server
	config        HTTPConfiguration
	serviceConfig ServiceConfiguration
}

func (srv *enhancedHTTPServer) listenAndServe() error {
	srv.server.Addr = srv.server.Addr

	if srv.config.CertFile != "" && srv.config.KeyFile != "" {
		return srv.server.ListenAndServeTLS(srv.config.Addr, srv.config.CertFile)
	}

	return srv.server.ListenAndServe()
}

func (srv *enhancedHTTPServer) StartWithContext(ctx context.Context) (err error) {
	done := make(chan bool, 1)

	go func() {
		<-ctx.Done()
		err = srv.server.Shutdown(context.Background())
		close(done)
	}()

	if err := srv.listenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-done
	return
}

func (srv *enhancedHTTPServer) Name() string {
	return srv.serviceConfig.Name
}

func (srv *enhancedHTTPServer) ApplyConfiguration(config interface{}) error {
	switch t := config.(type) {
	case *HTTPConfiguration:
		srv.config = *t
	case HTTPConfiguration:
		srv.config = t
	default:
		return rscsrv.ErrWrongConfigurationInformed
	}

	if srv.config.Addr == "" {
		return fmt.Errorf("addr not provided")
	}

	return nil
}

func (srv *enhancedHTTPServer) LoadConfiguration() (interface{}, error) {
	var target HTTPConfiguration
	if err := config.Load(srv.serviceConfig.Prefix, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// ServiceConfiguration TODO
type ServiceConfiguration struct {
	Name   string
	Prefix string
}

// EnhancedHTTPServer TODO
type EnhancedHTTPServer interface {
	rscsrv.Service
	rscsrv.Configurable
	rscsrv.StartableWithContext
}

// EnhanceHTTPServer enhances a http.Server into rscsrv.StartableWithContext.
func EnhanceHTTPServer(srv *http.Server, config ServiceConfiguration) EnhancedHTTPServer {
	if config.Prefix == "" {
		config.Prefix = "http"
	}

	if config.Name == "" {
		config.Name = "HTTP Server Service"
	}

	return &enhancedHTTPServer{
		server:        srv,
		serviceConfig: config,
	}
}
