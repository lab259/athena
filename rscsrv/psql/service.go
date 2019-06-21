package psqlrscsrv

import (
	"github.com/lab259/athena/config"
	psqlsrv "github.com/lab259/go-rscsrv-psql"
)

type PsqlService struct {
	psqlsrv.PsqlService
}

var DefaultPsqlService PsqlService

func (service *PsqlService) LoadConfiguration() (interface{}, error) {
	var configuration psqlsrv.Configuration
	err := config.Load("psql-service.yml", &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
