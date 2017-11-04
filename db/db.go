package db

import (
	"errors"
	"fmt"
	"github.com/nestorsokil/goto-url/util"
)

type DataSource interface {
	Find(key string) (*Record, error)
	FindShort(url string) (*Record, error)
	Save(*Record) error
	ExistsKey(key string) (bool, error)
	DeleteAllExpiredBefore(time int64) (removed int, err error)
	Update(*Record) error

	Shutdown()
}

type Record struct {
	Key        string
	URL        string
	Expiration int64
	MustExpire bool
}

func CreateDataSource(datasource string) (DataSource, error) {
	switch datasource {
	case util.IN_MEMORY:
		return NewMockDS()
	case util.MONGO:
		mongoConfig := util.LoadMongoConfig()
		return NewMongoDS(&mongoConfig)
	case util.REDIS:
		redisConfig := util.LoadRedisConfig()
		return NewRedisDs(&redisConfig)
	default:
		e := fmt.Sprintf("Unrecognized db option: %s", datasource)
		return nil, errors.New(e)
	}
}
