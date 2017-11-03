package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/util"
	log "github.com/sirupsen/logrus"
)

type UrlService struct {
	dataSource db.DataSource
	conf       *util.ApplicationConfig
}

// ErrNotFound is an error returned in case no suitable records were found
var ErrNotFound = errors.New("record not found")

// New returns a new UrlService
func New(dataSource db.DataSource, conf *util.ApplicationConfig) UrlService {
	return UrlService{dataSource, conf}
}

// RequestBuilder returns a new builder object for chained request creation
func (s *UrlService) RequestBuilder() *RequestBuilder {
	return builder(s.conf)
}

// GetRecord returns a record for the given request if there is one in data source
// If not - a new record is created
func (s *UrlService) GetRecord(request Request) (*db.Record, error) {
	if request.url == "" {
		return nil, errors.New("no URL provided")
	}
	var result *db.Record
	var err error
	if request.customKey != "" {
		result, err = s.createWithCustomKey(request.customKey, request.url, request.expire)
	} else {
		result, err = s.createWithRandKey(request.url, request.expire)
	}

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
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
	record, err := s.dataSource.FindShort(url)
	if record != nil {
		log.Debugf("URL '%s' already saved, responding with existing record", url)
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
		return nil, err
	}
	log.Debugf("Registered record with key '%s' for URL '%s'", key, url)
	return rec, nil
}

func getExpiration(expireIn int64) int64 {
	now := time.Now().UnixNano()
	expiration := now + expireIn*time.Hour.Nanoseconds()
	return expiration
}

// ClearRecordsAsync starts an infinite loop to check for expired records and delete them
// Use in separate goroutine to run in background
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
				log.Error(err.Error())
			} else if removed > 0 {
				log.Infof("Expired records removed. Count: %d", removed)
			}
		}
	}
}

// FindByKey returns a record for the provided key
func (s *UrlService) FindByKey(key string) (*db.Record, error) {
	record, err := s.dataSource.Find(key)
	if err != nil {
		log.Errorf("FindByKey(%s) : %v", key, err)
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
	log.Debugf("Record for URL '%s' requested.", record.URL)
	return record, nil
}

// ConstructURL creates a valid short URL
func (s *UrlService) ConstructURL(host, key string) string {
	var base string
	if s.conf.DevMode == true {
		base = s.conf.ApplicationUrl
	} else {
		base = host
	}
	return base + "/" + key
}
