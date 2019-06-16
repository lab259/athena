package rscsrv

import (
	"github.com/globalsign/mgo"

	"github.com/lab259/athena/config"
	mgosrv "github.com/lab259/go-rscsrv-mgo"
)

type MgoService struct {
	mgosrv.MgoService
}

func (service *MgoService) Name() string {
	return "Mgo Service"
}

var DefaultMgoService MgoService

func (service *MgoService) LoadConfiguration() (interface{}, error) {
	var configuration mgosrv.MgoServiceConfiguration
	err := config.Load("mgo-service.yml", &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

func (service *MgoService) RunWithDB(handler func(db *mgo.Database) error) error {
	return service.RunWithSession(func(session *mgo.Session) error {
		return handler(session.DB(service.Configuration.Database))
	})
}
