package service

import (
	"errors"
	"fmt"
	"github.com/nestorsokil/gl"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/util"
	"time"
)

type UrlService struct {
	dataSource db.DataSource
	conf       *util.ApplicationConfig
	logger     gl.Logger
}

func New(dataSource db.DataSource, conf *util.ApplicationConfig, logger gl.Logger) UrlService {
	return UrlService{dataSource, conf, logger}
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
	return s.createRecord(customKey, url, expireIn, true)
}

func (s *UrlService) createWithRandKey(url string, expireIn int64) (*db.Record, error) {
	record, _ := s.dataSource.FindShort(url)
	if record != nil {
		return record, nil
	}
	key, err := s.createKey()
	if err != nil {
		return nil, err
	}
	return s.createRecord(key, url, expireIn, false)
}

func (s *UrlService) createRecord(key, url string, expireIn int64, mustExpire bool) (*db.Record, error) {
	expiration := getExpiration(expireIn)
	rec := &db.Record{Key: key, URL: url, Expiration: expiration, MustExpire: mustExpire}
	err := s.dataSource.Save(rec)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return rec, nil
}

func getExpiration(expireIn int64) int64 {
	now := time.Now().UnixNano()
	expiration := now + expireIn*time.Hour.Nanoseconds()
	return expiration
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
			removed, err := s.dataSource.DeleteAllExpiredBefore(now)
			if err != nil {
				s.logger.Error(err.Error())
			} else if removed > 0 {
				s.logger.Info("Expired records removed. Count: %d", removed)
			}
		}
	}
}

var ErrNotFound = errors.New("Record not found.")

func (s *UrlService) FindByKey(key string) (*db.Record, error) {
	record, err := s.dataSource.Find(key)
	if err != nil {
		s.logger.Error("FindByKey(%s) : %v", key, err)
		return nil, ErrNotFound
	}
	if record == nil {
		return nil, ErrNotFound
	}
	if record.MustExpire {
		return record, nil
	}
	record.Expiration = getExpiration(s.conf.ExpirationTimeHours)
	s.dataSource.Update(record)
	return record, nil
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
