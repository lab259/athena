package psqltest

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/go-rscsrv"
	psqlsrv "github.com/lab259/go-rscsrv-psql"
	_ "github.com/lib/pq"
	"gopkg.in/src-d/go-kallax.v1"
)

func NewPsqlTestService() *PsqlTestService {
	var psqlTestService PsqlTestService

	service := psqlrscsrv.NewPsqlService()
	config, err := service.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	switch c := config.(type) {
	case psqlsrv.Configuration:
		psqlTestService.Configuration = c
	case *psqlsrv.Configuration:
		psqlTestService.Configuration = *c
	default:
		panic(rscsrv.ErrWrongConfigurationInformed)
	}

	psqlTestService.defaultService = service
	psqlTestService.id = kallax.NewULID()
	psqlTestService.identifier = fmt.Sprintf("txdb_%s", psqlTestService.id)
	txdb.Register(psqlTestService.identifier, "postgres", psqlTestService.Configuration.ConnectionString())

	return &psqlTestService
}

type PsqlTestService struct {
	started        bool
	id             kallax.ULID
	identifier     string
	defaultService psqlrscsrv.PsqlService
	db             *sql.DB
	Configuration  psqlsrv.Configuration
}

func (service *PsqlTestService) Name() string {
	return "Psql Test Service"
}

func (service *PsqlTestService) LoadConfiguration() (interface{}, error) {
	return service.defaultService.LoadConfiguration()
}

func (service *PsqlTestService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case psqlsrv.Configuration:
		service.Configuration = c
	case *psqlsrv.Configuration:
		service.Configuration = *c
	default:
		return rscsrv.ErrWrongConfigurationInformed
	}

	return nil

}

func (service *PsqlTestService) Restart() error {
	if service.db != nil {
		if err := service.Stop(); err != nil {
			return err
		}
	}
	return service.Start()
}

func (service *PsqlTestService) Start() error {
	if service.started {
		return nil
	}

	db, err := sql.Open(service.identifier, service.id.String())
	if err != nil {
		return err
	}

	if service.Configuration.MaxPoolSize > 0 {
		db.SetMaxOpenConns(service.Configuration.MaxPoolSize)
	}

	if err := db.Ping(); err != nil {
		return err
	}

	if err := psqlrscsrv.DefaultPsqlService.Stop(); err != nil {
		return err
	}

	service.db = db
	psqlrscsrv.DefaultPsqlService = service
	service.started = true
	return nil
}

func (service *PsqlTestService) Stop() error {
	if !service.started {
		return nil
	}

	if err := service.db.Close(); err != nil {
		return err
	}

	service.db = nil
	psqlrscsrv.DefaultPsqlService = service.defaultService
	return nil
}

func (service *PsqlTestService) Ping() error {
	if !service.started {
		return rscsrv.ErrServiceNotRunning
	}

	return service.db.Ping()
}

func (service *PsqlTestService) DB() (*sql.DB, error) {
	if !service.started {
		return nil, rscsrv.ErrServiceNotRunning
	}

	return service.db, nil
}
