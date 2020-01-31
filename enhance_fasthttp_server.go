package athena

import (
	"context"

	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/valyala/fasthttp"
)

type enhancedFasthttpServer struct {
	server *fasthttp.Server
	config FasthttpConfiguration
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
	return srv.config.Name
}

// FasthttpConfiguration holds additional configuration for the enhanced fasthttp.Server.
type FasthttpConfiguration struct {
	Name string
	Addr string
}

// EnhancedFasthttpServer TODO
type EnhancedFasthttpServer interface {
	rscsrv.Service
	rscsrv.StartableWithContext
}

// EnhanceFasthttpServer enhances a fasthttp.Server into rscsrv.StartableWithContext.
func EnhanceFasthttpServer(server *fasthttp.Server, config FasthttpConfiguration) EnhancedFasthttpServer {
	return &enhancedFasthttpServer{
		server: server,
		config: config,
	}
}
