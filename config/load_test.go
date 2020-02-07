package config_test

import (
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/lab259/athena/config"
	"github.com/lab259/athena/testing/envtest"
	"github.com/lab259/athena/testing/ginkgotest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLoad(t *testing.T) {
	ginkgotest.Init("athena/config Test Suite", t)
}

type ServiceConfiguration struct {
	Name    string        `yaml:"name"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

var _ = Describe("Load", func() {
	It("should load from YAML file (.yml)", func() {
		var srvConfig ServiceConfiguration
		Expect(config.Load("service.yml", &srvConfig)).To(Succeed())

		Expect(srvConfig.Name).To(Equal("Something"))
		Expect(srvConfig.Port).To(Equal(3303))
		Expect(srvConfig.Timeout).To(Equal(5 * time.Minute))
	})

	It("should load from YAML file (.yaml)", func() {
		var srvConfig ServiceConfiguration
		Expect(config.Load("service.yaml", &srvConfig)).To(Succeed())

		Expect(srvConfig.Name).To(Equal("Other"))
		Expect(srvConfig.Port).To(Equal(8080))
		Expect(srvConfig.Timeout).To(Equal(25 * time.Minute))
	})

	It("should override with environment variables", func() {
		envtest.With(map[string]string{
			"SERVICE_PORT":    "5656",
			"SERVICE_TIMEOUT": "5s",
		}, func() error {
			var srvConfig ServiceConfiguration
			Expect(config.Load("service.yml", &srvConfig)).To(Succeed())

			Expect(srvConfig.Name).To(Equal("Something"))
			Expect(srvConfig.Port).To(Equal(5656))
			Expect(srvConfig.Timeout).To(Equal(5 * time.Second))
			return nil
		})
	})

	It("should load from environment variables", func() {
		envtest.With(map[string]string{
			"SERVICE_NAME":    "Another",
			"SERVICE_PORT":    "9865",
			"SERVICE_TIMEOUT": "7s",
		}, func() error {
			var srvConfig ServiceConfiguration
			Expect(config.Load("service", &srvConfig)).To(Succeed())

			Expect(srvConfig.Name).To(Equal("Another"))
			Expect(srvConfig.Port).To(Equal(9865))
			Expect(srvConfig.Timeout).To(Equal(7 * time.Second))
			return nil
		})
	})
})

func projectFolder() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(file))
}
