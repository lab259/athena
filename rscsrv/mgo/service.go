package mgorscsrv

import (
	"github.com/globalsign/mgo"
	"github.com/lab259/go-rscsrv"

	"github.com/lab259/athena/config"
	mgosrv "github.com/lab259/go-rscsrv-mgo"
)

type MgoService interface {
	rscsrv.Service
	RunWithDB(handler func(db *mgo.Database) error) error
	RunWithSession(handler func(session *mgo.Session) error) error
}

type mgoService struct {
	mgosrv.MgoService
}

func (service *mgoService) Name() string {
	return "Mgo Service"
}

var DefaultMgoService MgoService = NewMgoService()

func (service *mgoService) LoadConfiguration() (interface{}, error) {
	var configuration mgosrv.MgoServiceConfiguration
	err := config.Load("mgo-service.yml", &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

func (service *mgoService) RunWithDB(handler func(db *mgo.Database) error) error {
	return service.RunWithSession(func(session *mgo.Session) error {
		return handler(session.DB(service.Configuration.Database))
	})
}

func NewMgoService() MgoService {
	return &mgoService{}
}
