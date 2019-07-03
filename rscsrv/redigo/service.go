package redigorscsrv

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/lab259/athena/config"
	"github.com/lab259/go-rscsrv"
	redigosrv "github.com/lab259/go-rscsrv-redigo"
)

type RedigoService interface {
	rscsrv.Service
	GetConn() (redis.Conn, error)
	RunWithConn(redigosrv.ConnHandler) error
	Publish(context.Context, string, interface{}) error
	Subscribe(context.Context, redigosrv.SubscribedHandler, redigosrv.SubscriptionHandler, ...string) error
}

var DefaultRedigoService RedigoService = NewRedigoService()

func NewRedigoService() RedigoService {
	return &redigoService{}
}

type redigoService struct {
	redigosrv.RedigoService
}

func (service *redigoService) Name() string {
	return "Redis Service"
}

func (service *redigoService) LoadConfiguration() (interface{}, error) {
	var configuration redigosrv.Configuration
	err := config.Load("redis-service.yml", &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
