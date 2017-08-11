package db

import (
	"github.com/fzzy/radix/redis"
	"github.com/nestorsokil/goto-url/util"
	"log"
	"strconv"
)

type RedisDataSource struct {
	client *redis.Client
}

func (rds *RedisDataSource) Find(key string) (*Record, error) {
	hash, err := rds.client.Cmd("HGETALL", "record:"+key).Hash()
	if err != nil {
		return nil, err
	}
	return asRecord(hash), nil
}

func (rds *RedisDataSource) FindShort(url string) (*Record, error) {
	hash, err := rds.client.Cmd("ZRANGEBYSCORE", "record.URL.index", url).Hash()
	if err != nil {
		return nil, err
	}
	return asRecord(hash), nil
}

func (rds *RedisDataSource) Save(newRecord *Record) error {
	key, url, exp := newRecord.Key, newRecord.URL, newRecord.Expiration
	reply := rds.client.Cmd("HMSET", "record:"+key, "key", key, "URL", url, "expiration", exp)
	if reply.Err != nil {
		return reply.Err
	}
	reply = rds.client.Cmd("ZADD", "record.URL.index", url, key)
	return nil
}

func (rds *RedisDataSource) ExistsKey(key string) (bool, error) {
	return rds.client.Cmd("EXISTS", "records:"+key).Bool()
}

func (rds *RedisDataSource) DeleteAllExpiredBefore(time int64) (removed int, err error) {
	// looks like redis is not a great choice for this
	replies := rds.client.Cmd("HGETALL", "record").Elems
	removed = 0
	for _, reply := range replies {
		hash, err := reply.Hash()
		if err != nil {
			continue
		}
		record := asRecord(hash)
		if record.Expiration < time {
			rds.client.Cmd("DEL", "record:"+record.Key)
			removed++
		}
	}
	return removed, nil
}

func (rds *RedisDataSource) Update(record *Record) error {
	return rds.Save(record)
}

func (rds *RedisDataSource) Shutdown() {
	rds.client.Close()
}

func asRecord(hash map[string]string) *Record {
	if len(hash) == 0 {
		return nil
	}
	record := new(Record)
	record.Key = hash["key"]
	record.URL = hash["URL"]
	exp, _ := strconv.Atoi(hash["expiration"])
	record.Expiration = int64(exp)
	return record
}

func NewRedisDs(config *util.RedisConfig) DataSource {
	var rds DataSource
	conn, err := redis.Dial("tcp", config.RedisUrl)
	if err != nil {
		log.Fatalf("Error creating Redis connection: %v", err)
	}

	rds = &RedisDataSource{conn}
	return rds
}