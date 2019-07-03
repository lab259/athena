package athena

import (
	"os"
	"os/signal"
	"syscall"

	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/valyala/fasthttp"

	"github.com/lab259/rlog"
)

func GracefulFastHTTP(srv *fasthttp.Server, addr string, serviceStarter rscsrv.ServiceStarter) {
	done := make(chan bool, 1)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		defer func() {
			serviceStarter.Stop(true)
			close(done)
		}()

		if err := srv.Shutdown(); err != nil {
			rlog.Criticalf("Failed closing idle connections: %v", err)
			os.Exit(1)
		}
	}()

	rlog.Infof("Starting binding the address %s", addr)
	if err := srv.ListenAndServe(addr); err != nil {
		rlog.Criticalf("Failed listing and serving: %v", err)
		os.Exit(1)
	}

	<-done
}
