package db

import (
	"github.com/nestorsokil/goto-url/util"
	"log"
)

type DataSource interface {
	// TODO refactor to (record, err)
	Find(key string) *Record
	// TODO refactor to (record, err)
	FindShort(url string) *Record
	Save(newRecord Record) error
	ExistsKey(key string) (bool, error)
	DeleteAllExpiredBefore(time int64) (removed int, err error)

	Shutdown()
}

type Record struct {
	Key        string
	URL        string
	Expiration int64
}

func CreateDataSource(config *util.ApplicationConfig) DataSource {
	dsType := config.Database
	switch dsType {
	case util.IN_MEMORY:
		return NewMockDS()
	case util.MONGO:
		mongoConfig := util.LoadMongoConfig()
		return NewMongoDS(&mongoConfig)
	case util.REDIS:
		redisConfig := util.LoadRedisConfig()
		return NewRedisDs(&redisConfig)
	default:
		log.Fatalf("Unrecognized db option: %s", dsType)
		return nil
	}
}
