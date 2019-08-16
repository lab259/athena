package rscsrvtest_test

import (
	"time"

	"github.com/lab259/athena/testing/rscsrvtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type clockService struct {
	startedAt time.Time
	stoppedAt time.Time
}

func (s *clockService) Name() string {
	panic("not implemented")
}

func (s *clockService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

func (s *clockService) Start() error {
	s.startedAt = time.Now()
	time.Sleep(time.Millisecond)
	return nil
}

func (s *clockService) Stop() error {
	s.stoppedAt = time.Now()
	time.Sleep(time.Millisecond)
	return nil
}

var _ = Describe("PsqlTestService", func() {
	var serv clockService
	var beforeCalled time.Time
	var afterCalled time.Time

	rscsrvtest.Before(func() {
		beforeCalled = time.Now()
		time.Sleep(time.Millisecond)
	}).After(func() {
		afterCalled = time.Now()
		time.Sleep(time.Millisecond)
	}).Setup(&serv)

	It("should run before hook after services are started", func() {
		Expect(beforeCalled).To(BeTemporally(">", serv.startedAt))
	})

	It("should run after hook before services are stopped", func() {
		Expect(afterCalled).To(BeTemporally("<", serv.stoppedAt))
	})
})
