package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nestorsokil/goto-url/conf"
	log "github.com/sirupsen/logrus"
)

type DataStorage interface {
	Find(ctx context.Context, key string) (*Record, error)
	SaveWithExpiration(ctx context.Context, record *Record, expireIn int64) error
	Exists(ctx context.Context, key string) (bool, error)

	Shutdown(ctx context.Context)
}

type Record struct {
	Key string
	URL string
}

func CreateStorage(c conf.Config) (DataStorage, error) {
	var (
		dataStorage DataStorage
		err         error
	)
	storage := c.GetString(conf.EnvStorage)
	switch storage {
	case conf.InMemory:
		dataStorage = newMockDS()
	case conf.Redis:
		dataStorage, err = newRedis(c)
		if err != nil {
			return nil, err
		}
	default:
		e := fmt.Sprintf("Unrecognized db option: %s", storage)
		return nil, errors.New(e)
	}

	if c.GetString(conf.EnvTraceDbEnabled) == "true" {
		return &traceDb{actual: dataStorage}, nil
	}
	return dataStorage, nil
}

func newMockDS() DataStorage {
	return &mockDataSource{records: make(map[string]*Record)}
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
