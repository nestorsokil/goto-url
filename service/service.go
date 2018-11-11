package service

import (
	"errors"
	"fmt"
	"github.com/nestorsokil/goto-url/conf"
	"github.com/nestorsokil/goto-url/db"
	log "github.com/sirupsen/logrus"
)

type UrlService struct {
	storage    db.DataStorage
	keyLen     int
	expiration int64
}

// ErrNotFound is an error returned in case no suitable records were found
var ErrNotFound = errors.New("record not found")

// New returns a new UrlService
func New(dataSource db.DataStorage, c conf.Config) UrlService {
	return UrlService{
		storage:    dataSource,
		keyLen:     c.GetInt(conf.ENV_KEY_LEN),
		expiration: int64(c.GetInt(conf.ENV_EXPIRATION_MILLIS)),
	}
}

// GetRecord returns a record for the given request if there is one in data source
// If not - a new record is created
func (s *UrlService) GetRecord(url, customKey string) (*db.Record, error) {
	if url == "" {
		return nil, errors.New("no URL provided")
	}
	var result *db.Record
	var err error
	if customKey != "" {
		result, err = s.createWithCustomKey(customKey, url)
	} else {
		result, err = s.createWithRandKey(url)
	}

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

// FindByKey returns a record for the provided key
func (s *UrlService) FindByKey(key string) (*db.Record, error) {
	record, err := s.storage.Find(key)
	if err != nil {
		log.Errorf("FindByKey(%s) : %v", key, err)
		return nil, ErrNotFound
	}
	if record == nil {
		return nil, ErrNotFound
	}
	s.storage.SaveWithExpiration(record, s.expiration)
	log.Debugf("Record for URL '%s' requested.", record.URL)
	return record, nil
}

// ConstructURL creates a valid short URL
func (s *UrlService) ConstructURL(host, key string) string {
	return host + "/" + key
}

func (s *UrlService) createWithCustomKey(customKey, url string) (*db.Record, error) {
	alreadyExists, err := s.storage.Exists(customKey)
	if err != nil {
		return nil, err
	}
	if alreadyExists {
		message := fmt.Sprintf("The custom key '%s' already exists", customKey)
		return nil, errors.New(message)
	}
	return s.createRecord(customKey, url, s.expiration)
}

func (s *UrlService) createWithRandKey(url string) (*db.Record, error) {
	key, err := s.createKey()
	if err != nil {
		return nil, err
	}
	return s.createRecord(key, url, s.expiration)
}

func (s *UrlService) createKey() (string, error) {
	var key string
	exists := true
	for exists {
		key = randKey(s.keyLen)
		e, err := s.storage.Exists(key)
		if err != nil {
			return "", err
		}
		exists = e
	}
	return key, nil
}

func (s *UrlService) createRecord(key, url string, expireIn int64) (*db.Record, error) {
	rec := &db.Record{Key: key, URL: url}
	err := s.storage.SaveWithExpiration(rec, expireIn)
	if err != nil {
		return nil, err
	}
	log.Debugf("Registered record with key '%s' for URL '%s'", key, url)
	return rec, nil
}
