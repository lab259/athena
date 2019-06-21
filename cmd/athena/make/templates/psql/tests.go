package templates_psql

import (
	"github.com/lab259/athena/cmd/athena/util/template"
)

var ServiceTestsTemplate = template.New("package_test.go", `package {{.Table}}_test

import (
	"testing"

	"github.com/lab259/athena/testing/ginkgotest"
)

func Test{{toCamel .Table}}(t *testing.T) {
	ginkgotest.Init("{{toCamel .Project}}/Services/{{toCamel .Table}} Test Suite", t)
}
`)
