package httprscsrv_test

import (
	"context"
	"testing"
	"time"

	"net/http"

	httprscsrv "github.com/lab259/athena/rscsrv/http"
	testingutils "github.com/lab259/athena/testing/ginkgotest"
	"github.com/lab259/athena/testing/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/cors"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestAthenaRscsrvHttp(t *testing.T) {
	testingutils.Init("Athena Rscsrv Http Test Suite", t)
}

var httpPort = ":3000"

func withHTTP(h func(*httpexpect.Expect)) func() {
	return httptest.WithHTTPReq("http://localhost"+httpPort, h)
}

var _ = Describe("Http", func() {
	It("should start an http server", func(done Done) {
		withHTTP(func(app *httpexpect.Expect) {
			mux := http.NewServeMux()

			mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("done!"))
			})

			httpServer := httprscsrv.NewHttpServer(&http.Server{
				Addr:    httpPort,
				Handler: cors.AllowAll().Handler(mux),
			}, httprscsrv.ServiceConfiguration{
				Name: "Admin",
			})

			config, err := httpServer.LoadConfiguration()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(httpServer.ApplyConfiguration(config)).ShouldNot(HaveOccurred())
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				defer GinkgoRecover()
				err := httpServer.StartWithContext(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(http.ErrServerClosed))

				close(done) // Tells the ginkgo the test is done.
			}()

			time.Sleep(10 * time.Millisecond)

			app.GET("/test").Expect().Status(http.StatusOK).Body().Equal("done!")

			cancel()
		})()
	}, 1)
})
