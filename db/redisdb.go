package db

import (
	"errors"
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/nestorsokil/goto-url/util"
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
	hash, err := rds.client.Cmd("HGETALL", "record.URL:"+url).Hash()
	if err != nil {
		return nil, err
	}
	return rds.Find(hash["key"])
}

func (rds *RedisDataSource) Save(newRecord *Record) error {
	key, url, exp, mustExp := newRecord.Key, newRecord.URL,
		newRecord.Expiration, newRecord.MustExpire
	reply := rds.client.Cmd(
		"HMSET", "record:"+key, "key", key, "URL", url,
		"expiration", exp, "mustExpire", mustExp)
	if reply.Err != nil {
		return reply.Err
	}
	reply = rds.client.Cmd("HMSET", "record.URL:"+url, "key", key)
	if reply.Err != nil {
		return reply.Err
	}
	reply = rds.client.Cmd("ZADD", "record.time.index", exp, key)
	return reply.Err
}

func (rds *RedisDataSource) ExistsKey(key string) (bool, error) {
	return rds.client.Cmd("EXISTS", "record:"+key).Bool()
}

func (rds *RedisDataSource) DeleteAllExpiredBefore(time int64) (removed int, err error) {
	replies := rds.client.Cmd("ZRANGEBYSCORE", "record.time.index", 0, time).Elems
	removed = 0
	for _, reply := range replies {
		if reply == nil {
			continue
		}
		hash, err := reply.Hash()
		if err != nil {
			continue
		}
		record, err := rds.Find(hash["key"])
		if err != nil {
			continue
		}
		rds.client.Cmd("DEL", "record:"+record.Key)
		rds.client.Cmd("DEL", "record.URL:"+record.URL)
		rds.client.Cmd("ZREM", "record.time.index", record.Key)
		removed++
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

func NewRedisDs(config *util.RedisConfig) (DataSource, error) {
	conn, err := redis.Dial("tcp", config.RedisUrl)
	if err != nil {
		e := fmt.Sprintf("Error creating Redis connection: %v", err)
		return nil, errors.New(e)
	}
	return &RedisDataSource{conn}, nil
}
