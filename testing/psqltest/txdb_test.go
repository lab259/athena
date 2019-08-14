package psqltest_test

import (
	"testing"

	"github.com/lab259/athena/testing/ginkgotest"
	"github.com/lab259/athena/testing/psqltest"
	"github.com/lab259/go-rscsrv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPsqlTestService(t *testing.T) {
	ginkgotest.Init("athena/testing/psqltest Test Suite", t)
}

var _ = Describe("PsqlTestService", func() {
	It("should load and apply configuration", func() {
		psqlTestService := psqltest.NewPsqlTestService()
		config, err := psqlTestService.LoadConfiguration()
		Expect(err).ToNot(HaveOccurred())
		Expect(psqlTestService.ApplyConfiguration(config)).To(Succeed())
	})

	It("should start and stop configuration", func() {
		psqlTestService := psqltest.NewPsqlTestService()
		config, err := psqlTestService.LoadConfiguration()
		Expect(err).ToNot(HaveOccurred())
		Expect(psqlTestService.ApplyConfiguration(config)).To(Succeed())
		Expect(psqlTestService.Start()).To(Succeed())
		Expect(psqlTestService.Ping()).To(Succeed())
		Expect(psqlTestService.Stop()).To(Succeed())
		Expect(psqlTestService.Ping()).To(Equal(rscsrv.ErrServiceNotRunning))
	})
})
