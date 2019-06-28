package httptest

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"gopkg.in/gavv/httpexpect.v1"
)

type Handler func(*httpexpect.Expect)

func WithFastHTTP(handler fasthttp.RequestHandler, cb Handler) func() {
	return func() {
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(handler),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: &httpGomegaFail{2},
		})

		cb(expect)
	}
}
