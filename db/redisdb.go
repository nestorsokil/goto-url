package db

import (
	"github.com/fzzy/radix/redis"
	log "github.com/sirupsen/logrus"
)

type redisdb struct {
	client *redis.Client
}

func (r *redisdb) Find(key string) (*Record, error) {
	hash, err := r.client.Cmd("HGETALL", key).Hash()
	if err != nil {
		log.Errorf("Error when looking up record (key = %v)", key)
		return nil, err
	}
	if len(hash) == 0 {
		return nil, nil
	}

	return &Record{Key: key, URL: hash["URL"]}, nil
}

func (r *redisdb) SaveWithExpiration(rec *Record, expireIn int64) error {
	reply := r.client.Cmd("HMSET", rec.Key, "URL", rec.URL)
	if reply.Err != nil {
		log.Errorf("Error when saving record (record = %v, expireIn = %v)", rec, expireIn)
		return reply.Err
	}
	return nil
}

func (r *redisdb) Exists(key string) (bool, error) {
	return r.client.Cmd("EXISTS", key).Bool()
}

func (r *redisdb) Shutdown() {
	r.client.Close()
}
