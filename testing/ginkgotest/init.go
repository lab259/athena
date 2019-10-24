package ginkgotest

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/athena/config"
	"github.com/lab259/athena/testing/envtest"
	"github.com/lab259/rlog/v2"
	"github.com/onsi/ginkgo"
	ginkgoConfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

type SetWriter func(io.Writer)

func Init(description string, t *testing.T, loggers ...SetWriter) {
	log.SetOutput(ginkgo.GinkgoWriter)
	rlog.SetOutput(ginkgo.GinkgoWriter)
	for _, logger := range loggers {
		logger(ginkgo.GinkgoWriter)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}

	envtest.Override(map[string]string{
		"PROJECT_ROOT": config.ProjectRoot(),
		"ENV":          env,
	}, func() error {
		dir, _ := os.Getwd()
		ginkgo.GinkgoWriter.Write([]byte(fmt.Sprintf("CWD: %s\n", dir)))
		ginkgo.GinkgoWriter.Write([]byte(fmt.Sprintf("ENV: %s\n", os.Getenv("ENV"))))
		ginkgo.GinkgoWriter.Write([]byte(fmt.Sprintf("Random Seed: %d\n", ginkgoConfig.GinkgoConfig.RandomSeed)))
		gomega.RegisterFailHandler(ginkgo.Fail)

		if os.Getenv("CI") == "" {
			macchiato.RunSpecs(t, description)
		} else {
			projectRoot := config.ProjectRoot()
			project := filepath.Base(projectRoot)
			reporterOutputDir := path.Join(projectRoot, "test-results", project, strings.Replace(dir, projectRoot, "", 1))
			os.MkdirAll(reporterOutputDir, os.ModePerm)
			junitReporter := reporters.NewJUnitReporter(path.Join(reporterOutputDir, "results.xml"))
			macchiatoReporter := macchiato.NewReporter()
			ginkgo.RunSpecsWithCustomReporters(t, description, []ginkgo.Reporter{
				macchiatoReporter,
				junitReporter,
			})
		}
		return nil
	})
}
