package httptest

import (
	"gopkg.in/gavv/httpexpect.v1"
)

func WithHTTPReq(baseURL string, cb func(*httpexpect.Expect)) func() {
	return func() {
		expect := httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  baseURL,
			Reporter: &httpGomegaFail{2},
		})
		cb(expect)
	}
}
