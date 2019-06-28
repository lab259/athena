package httptest

import (
	"net/http"

	"github.com/lab259/hermes"
	"gopkg.in/gavv/httpexpect.v1"
)

func WithHermes(router hermes.Router, cb Handler) func() {
	return func() {
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(router.Handler()),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: &httpGomegaFail{2},
		})

		cb(expect)
	}
}
