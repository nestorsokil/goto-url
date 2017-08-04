package service

import (
	"time"
	"log"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/util"
)

type UrlService struct {
	dataSource db.DataSource
	conf *util.Configuration
}

func New(dataSource db.DataSource, conf *util.Configuration) UrlService {
	return UrlService{dataSource, conf}
}

func (s *UrlService) GetRecord(url string) (*db.Record, error) {
	record := s.dataSource.FindShort(url)
	if record != nil {
		return record, nil
	}
	// persist new record
	var key string
	exists := true
	for exists {
		key = randKey(s.conf.KeyLength)
		e, err := s.dataSource.ExistsKey(key)
		if err != nil {
			return nil, err
		}
		exists = e
	}
	now := time.Now().UnixNano()
	expiresIn := s.conf.ExpirationTimeHours * time.Hour.Nanoseconds()
	expiration := now + expiresIn
	rec := db.Record{Key:key, URL:url, Expiration:expiration}
	err := s.dataSource.Save(rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *UrlService) ClearRecordsAsync(stopSignal <-chan struct{}) {
	waitTime := time.Duration(s.conf.ClearTimeSeconds * time.Second.Nanoseconds())
	for {
		select {
		case <-stopSignal:
			return
		default:
			time.Sleep(waitTime)
			now := time.Now().UnixNano()
			removed, err := s.dataSource.DeleteAllAfter(now)
			if err != nil {
				log.Println("[ERROR]", err)
			}
			log.Println("[INFO] Expired records removed. Count:", removed)
		}
	}
}

func (s *UrlService) FindByKey(key string) *db.Record {
	return s.dataSource.Find(key)
}