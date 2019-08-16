package rscsrvtest

import (
	"github.com/lab259/go-rscsrv"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

type Setupper interface {
	Setup(...rscsrv.Service)
	Before(func()) Setupper
	After(func()) Setupper
}

type starter struct {
	offset int
	before []func()
	after  []func()
}

func (s *starter) Before(fn func()) Setupper {
	if s.before == nil {
		s.before = make([]func(), 0, 5)
	}
	s.before = append(s.before, fn)
	return s
}

func (s *starter) After(fn func()) Setupper {
	if s.after == nil {
		s.after = make([]func(), 0, 5)
	}
	s.after = append(s.after, fn)
	return s
}

func (s *starter) Setup(services ...rscsrv.Service) {
	serviceStarter := rscsrv.NewServiceStarter(&rscsrv.NopStarterReporter{}, services...)

	ginkgo.BeforeEach(func() {
		gomega.ExpectWithOffset(s.offset+1, serviceStarter.Start()).To(gomega.Succeed())

		for _, beforeFn := range s.before {
			beforeFn()
		}
	})

	ginkgo.AfterEach(func() {
		for _, afterFn := range s.after {
			afterFn()
		}

		gomega.ExpectWithOffset(s.offset+1, serviceStarter.Stop(false)).To(gomega.Succeed())
	})
}

func Setup(services ...rscsrv.Service) {
	s := starter{
		offset: 1,
	}
	s.Setup(services...)
}

func Before(fn func()) Setupper {
	s := new(starter)
	s.before = []func(){fn}
	return s
}

func After(fn func()) Setupper {
	s := new(starter)
	s.after = []func(){fn}
	return s
}
