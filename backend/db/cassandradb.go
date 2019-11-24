package db

import (
	"context"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type cassandraDb struct {
	session *gocql.Session
}

func (c *cassandraDb) Find(ctx context.Context, key string) (*Record, error) {
	var r Record
	query := c.session.Query(`SELECT (key, url) WHERE key = ?`, key)
	if err := query.Scan(&r.Key, &r.URL); err != nil {
		log.Errorf("Error looking up row: %v", err.Error())
		return nil, err
	}
	return &r, nil
}

func (c *cassandraDb) SaveWithExpiration(ctx context.Context, record *Record, expireIn int64) error {
	ttlSeconds := expireIn / 1000
	query := c.session.Query(`INSERT INTO records (key, url) VALUES (?, ?) USING TTL ?;`,
		record.Key, record.URL, ttlSeconds)
	if err := query.Exec(); err != nil {
		log.Errorf("Error inserting row: %v", err.Error())
		return err
	}
	return nil
}

func (c *cassandraDb) Exists(ctx context.Context, key string) (bool, error) {
	var exists int
	query := c.session.Query(`SELECT COUNT(*) WHERE key = ?`, key)
	if err := query.Scan(&exists); err != nil {
		log.Error("Error checking for existence: %v", err.Error())
		return false, err
	}
	return exists > 0, nil
}

func (c *cassandraDb) Shutdown(ctx context.Context) {
	panic("implement me")
}
