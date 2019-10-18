package sqltest_test

import (
	"testing"

	"github.com/lab259/athena/testing/ginkgotest"
)

func TestSqltest(t *testing.T) {
	ginkgotest.Init("athena/testing/sqltest Test Suite", t)
}
