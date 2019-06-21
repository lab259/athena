package psqltest

import (
	"database/sql"
	"fmt"

	psqlrscsrv "github.com/lab259/athena/rscsrv/psql"
	"github.com/lab259/go-rscsrv"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/DATA-DOG/go-txdb"
	psqlsrv "github.com/lab259/go-rscsrv-psql"
	_ "github.com/lib/pq"
)

func NewPsqlTestService() *PsqlTestService {
	return &PsqlTestService{}
}

type PsqlTestService struct {
	db            *sql.DB
	Configuration psqlsrv.Configuration
}

func (service *PsqlTestService) Name() string {
	return "Psql Test Service"
}

func (service *PsqlTestService) LoadConfiguration() (interface{}, error) {
	return psqlrscsrv.NewPsqlService().LoadConfiguration()
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
	if service.db != nil {
		if err := service.Stop(); err != nil {
			return err
		}
	}

	ulid := kallax.NewULID()
	identifier := fmt.Sprintf("txdb_%s", ulid)
	txdb.Register(identifier, "postgres", service.Configuration.ConnectionString())

	db, err := sql.Open(identifier, ulid.String())
	if err != nil {
		return err
	}

	if service.Configuration.MaxPoolSize > 0 {
		db.SetMaxOpenConns(service.Configuration.MaxPoolSize)
	}

	if err := db.Ping(); err != nil {
		return err
	}

	service.db = db
	psqlrscsrv.DefaultPsqlService = service
	return nil
}

func (service *PsqlTestService) Stop() error {
	if service.db == nil {
		return nil
	}

	if err := service.db.Close(); err != nil {
		return err
	}

	service.db = nil
	psqlrscsrv.DefaultPsqlService = psqlrscsrv.NewPsqlService()
	return nil
}

func (service *PsqlTestService) Ping() error {
	if service.db == nil {
		return rscsrv.ErrServiceNotRunning
	}

	return service.db.Ping()
}

func (service *PsqlTestService) DB() (*sql.DB, error) {
	if service.db == nil {
		return nil, rscsrv.ErrServiceNotRunning
	}

	return service.db, nil
}
