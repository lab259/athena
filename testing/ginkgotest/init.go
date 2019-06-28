package ginkgotest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/athena/config"
	"github.com/lab259/athena/testing/envtest"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

func Init(description string, t *testing.T) {
	envtest.Override(map[string]string{
		"PROJECT_ROOT": config.ProjectRoot(),
	}, func() error {
		if os.Getenv("ENV") == "" {
			err := os.Setenv("ENV", "test")
			if err != nil {
				panic(err)
			}
		}

		dir, _ := os.Getwd()
		ginkgo.GinkgoWriter.Write([]byte(fmt.Sprintf("CWD: %s\n", dir)))
		ginkgo.GinkgoWriter.Write([]byte(fmt.Sprintf("ENV: %s\n", os.Getenv("ENV"))))
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
