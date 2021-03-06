package httptest

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/onsi/ginkgo"
)

func pruneStack(fullStackTrace string, skip int) int {
	stack := strings.Split(fullStackTrace, "\n")
	if len(stack) > 2*(skip+1) {
		stack = stack[2*(skip+1):]
	}
	prunedStack := []string{}
	re := regexp.MustCompile(`\/ginkgo\/|\/pkg\/testing\/|\/pkg\/runtime\/|\/gavv\/httpexpect`)
	for i := 0; i < len(stack)/2; i++ {
		if !re.Match([]byte(stack[i*2])) && i >= skip {
			return i
		}
		prunedStack = append(prunedStack, stack[i*2])
		prunedStack = append(prunedStack, stack[i*2+1])
	}
	return -1
}

type httpGomegaFail struct {
	Skip int
}

func (r *httpGomegaFail) Errorf(message string, args ...interface{}) {
	ginkgo.Fail(fmt.Sprintf(message, args...), pruneStack(string(debug.Stack()), r.Skip)+1)
}
