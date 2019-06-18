package templates_mongo

import (
	"github.com/lab259/athena/athena/util/template"
)

var ServiceTestsTemplate = template.New("package_test.go", `package {{.Collection}}_test

import (
	"testing"

	"github.com/lab259/athena/testing/ginkgo"
)

func Test{{toCamel .Collection}}(t *testing.T) {
	ginkgo.Init("{{toCamel .Project}}/Services/{{toCamel .Collection}} Test Suite", t)
}
`)
