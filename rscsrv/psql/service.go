package psqlrscsrv

import (
	"database/sql"

	"github.com/lab259/athena/config"
	"github.com/lab259/go-rscsrv"
	psqlsrv "github.com/lab259/go-rscsrv-psql"
)

type PsqlService interface {
	rscsrv.Service
	rscsrv.Configurable
	rscsrv.Startable
	rscsrv.Stoppable
	DB() (*sql.DB, error)
	Ping() error
}

type psqlService struct {
	psqlsrv.PsqlService
}

var DefaultPsqlService PsqlService = NewPsqlService()

func (service *psqlService) LoadConfiguration() (interface{}, error) {
	var configuration psqlsrv.Configuration
	err := config.Load("psql-service.yml", &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

func NewPsqlService() PsqlService {
	return &psqlService{}
}
