package db

import (
	"github.com/nestorsokil/goto-url/util"
	"log"
)

type DataSource interface {
	Find(key string) *Record
	FindShort(url string) *Record
	Save(newRecord Record) error
	ExistsKey(key string) (bool, error)
	DeleteAllAfter(time int64) (removed int, err error)

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
	default:
		log.Fatalf("Unrecognized db option: %s", dsType)
		return nil
	}
}
