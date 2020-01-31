package athena

import (
	"context"
	"net/http"

	rscsrv "github.com/lab259/go-rscsrv"
)

type enhancedHTTPServer struct {
	server *http.Server
}

func (srv *enhancedHTTPServer) StartWithContext(ctx context.Context) (err error) {
	done := make(chan bool, 1)

	go func() {
		<-ctx.Done()
		err = srv.server.Shutdown(context.Background())
		close(done)
	}()

	if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-done
	return
}

// EnhanceHTTPServer enhances a http.Server into rscsrv.StartableWithContext.
func EnhanceHTTPServer(srv *http.Server) rscsrv.StartableWithContext {
	return &enhancedHTTPServer{
		server: srv,
	}
}
