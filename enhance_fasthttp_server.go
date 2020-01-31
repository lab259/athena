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

// FasthttpConfiguration holds additional configuration for the enhanced fasthttp.Server.
type FasthttpConfiguration struct {
	Addr string
}

// EnhanceFasthttpServer enhances a fasthttp.Server into rscsrv.StartableWithContext.
func EnhanceFasthttpServer(server *fasthttp.Server, config FasthttpConfiguration) rscsrv.StartableWithContext {
	return &enhancedFasthttpServer{
		server: server,
		config: config,
	}
}
