package service

import (
	"errors"
	"fmt"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/util"
	"log"
	"time"
)

type UrlService struct {
	dataSource db.DataSource
	conf       *util.Configuration
}

func New(dataSource db.DataSource, conf *util.Configuration) UrlService {
	return UrlService{dataSource, conf}
}

func (s *UrlService) RequestBuilder() *RequestBuilder {
	return builder(s.conf)
}

func (s *UrlService) GetRecord(request Request) (*db.Record, error) {
	if request.url == "" {
		return nil, errors.New("No URL provided.")
	}
	if request.customKey != "" {
		return s.createWithCustomKey(request.customKey, request.url, request.expire)
	}
	return s.createWithRandKey(request.url, request.expire)
}

func (s *UrlService) createKey() (string, error) {
	var key string
	exists := true
	for exists {
		key = util.RandKey(s.conf.KeyLength)
		e, err := s.dataSource.ExistsKey(key)
		if err != nil {
			return "", err
		}
		exists = e
	}
	return key, nil
}

func (s *UrlService) createWithCustomKey(customKey, url string, expireIn int64) (*db.Record, error) {
	alreadyExists, err := s.dataSource.ExistsKey(customKey)
	if err != nil {
		return nil, err
	}
	if alreadyExists {
		message := fmt.Sprintf("The custom key '%s' already exists", customKey)
		return nil, errors.New(message)
	}
	return s.createRecord(customKey, url, expireIn)
}

func (s *UrlService) createWithRandKey(url string, expireIn int64) (*db.Record, error) {
	record := s.dataSource.FindShort(url)
	if record != nil {
		return record, nil
	}
	key, err := s.createKey()
	if err != nil {
		return nil, err
	}
	return s.createRecord(key, url, expireIn)
}

func (s *UrlService) createRecord(key, url string, expireIn int64) (*db.Record, error) {
	now := time.Now().UnixNano()
	expiration := now + expireIn*time.Hour.Nanoseconds()
	rec := db.Record{Key: key, URL: url, Expiration: expiration}
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

func (s *UrlService) ConstructUrl(host, key string) string {
	var base string
	if s.conf.DevMode == true {
		base = s.conf.ApplicationUrl
	} else {
		base = host
	}
	return base + "/" + key
}
