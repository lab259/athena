package rscsrvtest

import (
	"github.com/lab259/go-rscsrv"
	"github.com/onsi/gomega"
)

func Start(services ...rscsrv.Service) {
	gomega.ExpectWithOffset(1, rscsrv.NewServiceStarter(services, &rscsrv.NopServiceReporter{}).Start()).To(gomega.Succeed())
}
