package athena

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	rscsrv "github.com/lab259/go-rscsrv"

	"github.com/lab259/rlog/v2"
)

func GracefulHTTP(srv *http.Server, serviceStarter rscsrv.ServiceStarter) {
	done := make(chan bool, 1)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		defer func() {
			serviceStarter.Stop(true)
			close(done)
		}()

		if err := srv.Shutdown(context.Background()); err != nil {
			rlog.Criticalf("Failed closing idle connections: %v", err)
			os.Exit(1)
		}
	}()

	rlog.Infof("Starting binding the address %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		rlog.Criticalf("Failed listing and serving: %v", err)
		os.Exit(1)
	}

	<-done
}
