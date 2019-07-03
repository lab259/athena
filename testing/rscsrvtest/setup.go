package rscsrvtest

import (
	"github.com/lab259/go-rscsrv"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func Setup(services ...rscsrv.Service) {
	serviceStarter := rscsrv.NewServiceStarter(&rscsrv.NopStarterReporter{}, services...)

	ginkgo.BeforeEach(func() {
		gomega.ExpectWithOffset(1, serviceStarter.Start()).To(gomega.Succeed())
	})

	ginkgo.AfterEach(func() {
		gomega.ExpectWithOffset(1, serviceStarter.Stop(false)).To(gomega.Succeed())
	})
}
