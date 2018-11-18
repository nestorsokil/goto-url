package db

import (
	"errors"
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/nestorsokil/goto-url/conf"
	log "github.com/sirupsen/logrus"
)

type DataStorage interface {
	Find(key string) (*Record, error)
	SaveWithExpiration(record *Record, expireIn int64) error
	Exists(key string) (bool, error)

	Shutdown()
}

type Record struct {
	Key string
	URL string
}

func CreateStorage(c conf.Config) (DataStorage, error) {
	storage := c.GetString(conf.ENV_STORAGE)
	switch storage {
	case conf.IN_MEMORY:
		return newMockDS()
	case conf.REDIS:
		return newRedis(c.GetString(conf.ENV_REDIS_URL))
	default:
		e := fmt.Sprintf("Unrecognized db option: %s", storage)
		return nil, errors.New(e)
	}
}

func newMockDS() (DataStorage, error) {
	return &mockDataSource{records: make(map[string]*Record)}, nil
}

func newRedis(address string) (DataStorage, error) {
	log.Infof("Trying to connect to Redis on address '%v'", address)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		e := fmt.Sprintf("Error creating Redis connection: %v", err)
		return nil, errors.New(e)
	}
	return &redisdb{conn}, nil
}
