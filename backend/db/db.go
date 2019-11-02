package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nestorsokil/goto-url/conf"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
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

	if strings.ToLower(c.GetString(conf.EnvTraceDbEnabled)) == "true" {
		return &traceDb{actual: dataStorage}, nil
	}
	return dataStorage, nil
}

func newMockDS() DataStorage {
	return &mockDataSource{records: make(map[string]*Record)}
}

const attempts = 5

func newRedis(c conf.Config) (DataStorage, error) {
	address, password := c.GetString(conf.EnvRedisUrl), c.GetString(conf.EnvRedisPass)
	log.Infof("Trying to connect to Redis on address '%v'", address)
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
	for i := 1; i <= attempts; i++ {
		pong, err := cli.Ping().Result()
		if err == nil {
			return &redisdb{cli}, nil
		}
		log.Errorf("Could not create Redis client (attempt = %v, response = %v, error = %v)", i, pong, err)
		time.Sleep(5 * time.Second) // todo exponential
	}
	log.Errorf("Could not establish Redis connection in %v attempts", attempts)
	return nil, errors.New("connection error")
}
