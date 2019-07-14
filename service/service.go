package service

import (
	"errors"
	"fmt"
	"github.com/nestorsokil/goto-url/conf"
	"github.com/nestorsokil/goto-url/db"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type UrlService struct {
	storage    db.DataStorage
	keyLen     int
	expiration int64
}

var ErrNotFound = errors.New("record not found")
var ErrNotUrl = errors.New("not valid URL")

var urlRegex = regexp.MustCompile("(https://|http://)?(?:([\\w-]+)\\.)?([\\w-]+)\\.(\\w+)")

// New returns a new UrlService
func New(dataSource db.DataStorage, c conf.Config) UrlService {
	return UrlService{
		storage:    dataSource,
		keyLen:     c.GetInt(conf.EnvKeyLen),
		expiration: int64(c.GetInt(conf.EnvExpirationMillis)),
	}
}

// CreateRecord creates a record for the given url and (optional) key pair
func (s *UrlService) CreateRecord(rawUrl, customKey string) (*db.Record, error) {
	if rawUrl == "" {
		return nil, errors.New("no URL provided")
	}
	url, ok := validateUrl(rawUrl)
	if !ok {
		log.Errorf("The provided string '%v' is not a url", rawUrl)
		return nil, ErrNotUrl
	}
	var result *db.Record
	var err error
	if customKey != "" {
		result, err = s.createWithCustomKey(customKey, url)
	} else {
		result, err = s.createWithRandKey(url)
	}
	if err != nil {
		log.Errorf("Error saving the record (url = %v, key = %v): %v", url, customKey, err.Error())
		return nil, err
	}
	log.Debugf("Record %v was created successfully", result)
	return result, nil
}

func validateUrl(url string) (validated string, isCorrectUrl bool) {
	matches := urlRegex.FindAllStringSubmatch(url, -1)
	if matches == nil || len(matches) < 1 {
		return "", false
	}
	if matches[0][1] == "" {
		return "http://" + url, true
	}
	return url, true
}

// FindByKey returns a record for the provided key
func (s *UrlService) FindByKey(key string) (*db.Record, error) {
	record, err := s.storage.Find(key)
	if err != nil {
		log.Errorf("FindByKey(%s) : %v", key, err)
		return nil, ErrNotFound
	}
	if record == nil {
		log.Errorf("Key '%v' not found", key)
		return nil, ErrNotFound
	}
	if err = s.storage.SaveWithExpiration(record, s.expiration); err != nil {
		log.Warnf("Could not refresh record (key = %s)", key)
	}
	log.Debugf("Record for URL '%s' requested.", record.URL)
	return record, nil
}

func (s *UrlService) createWithCustomKey(customKey, url string) (*db.Record, error) {
	alreadyExists, err := s.storage.Exists(customKey)
	if err != nil {
		log.Errorf("Error check the db for key '%v', URL '%v'. Error: %v", customKey, url, err)
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
	log.Debugf("Random key '%v' was generated for URL '%v'", key, url)
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
	log.Debugf("Creating record (key = %v, URL = %v)", key, url)
	rec := &db.Record{Key: key, URL: url}
	err := s.storage.SaveWithExpiration(rec, expireIn)
	if err != nil {
		log.Errorf("Could not save record %v. Error: %v", rec, err)
		return nil, err
	}
	log.Debugf("Registered record with key '%s' for URL '%s'", key, url)
	return rec, nil
}
