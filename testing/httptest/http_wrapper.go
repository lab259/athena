package httptest

import (
	"net/http"

	"gopkg.in/gavv/httpexpect.v1"
)

func WithHTTP(handler http.Handler, cb Handler) func() {
	return func() {
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: &httpGomegaFail{2},
		})

		cb(expect)
	}
}
