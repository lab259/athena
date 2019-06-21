package templates_mongo

import (
	"github.com/lab259/athena/cmd/athena/util/template"
)

var ServiceTestsTemplate = template.New("package_test.go", `package {{.Collection}}_test

import (
	"testing"

	"github.com/lab259/athena/testing/ginkgotest"
)

func Test{{toCamel .Collection}}(t *testing.T) {
	ginkgotest.Init("{{toCamel .Project}}/Services/{{toCamel .Collection}} Test Suite", t)
}
`)
