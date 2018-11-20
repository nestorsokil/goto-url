package db

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
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
	storage := c.GetString(conf.EnvStorage)
	switch storage {
	case conf.InMemory:
		return newMockDS()
	case conf.Redis:
		return newRedis(c)
	default:
		e := fmt.Sprintf("Unrecognized db option: %s", storage)
		return nil, errors.New(e)
	}
}

func newMockDS() (DataStorage, error) {
	return &mockDataSource{records: make(map[string]*Record)}, nil
}

func newRedis(c conf.Config) (DataStorage, error) {
	address, password := c.GetString(conf.EnvRedisUrl), c.GetString(conf.EnvRedisPass)
	log.Infof("Trying to connect to Redis on address '%v'", address)
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
	pong, err := cli.Ping().Result()
	if err != nil {
		log.Errorf("Could not create Redis client (response = %v, error = %v)", pong, err)
		return nil, errors.New("connection error")
	}
	return &redisdb{cli}, nil
}
