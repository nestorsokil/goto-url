package db

import (
	"context"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type redisdb struct {
	client *redis.Client
}

func (r *redisdb) Find(ctx context.Context, key string) (*Record, error) {
	hash, err := r.client.HGetAll(key).Result()
	if err != nil {
		log.Errorf("Error when looking up record (key = %v)", key)
		return nil, err
	}
	if len(hash) == 0 {
		return nil, nil
	}
	return &Record{Key: key, URL: hash["URL"]}, nil
}

func (r *redisdb) SaveWithExpiration(ctx context.Context, rec *Record, expireIn int64) error {
	params := map[string]interface{}{"URL": rec.URL}
	return r.client.HMSet(rec.Key, params).Err()
}

func (r *redisdb) Exists(ctx context.Context, key string) (bool, error) {
	i, err := r.client.Exists(key).Result()
	return i == 1, err
}

func (r *redisdb) Shutdown(ctx context.Context) {
	err := r.client.Close()
	if err != nil {
		log.Errorf("Could not close the client properly %v", err)
	}
}
