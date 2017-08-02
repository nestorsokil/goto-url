package service

import (
	"time"
	"log"
	"math/rand"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/config"
)

var src = rand.NewSource(time.Now().UnixNano())
const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6
	letterIdxMask = 1 << letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

type Service struct {
	dataSource db.DataSource
}

func New(dataSource db.DataSource) Service {
	return Service{dataSource}
}

func (s *Service) GetRecord(url string) (*db.Record, error) {
	record := s.dataSource.FindShort(url)
	if record != nil {
		return record, nil
	}
	// persist new record
	var key string
	exists := true
	for exists {
		key = s.randKey()
		e, err := s.dataSource.ExistsKey(key)
		if err != nil {
			return nil, err
		}
		exists = e
	}
	expiresIn := time.Duration(config.Settings.ExpirationTimeHours) * time.Hour
	expiration := uint64(time.Now().Add(expiresIn).Unix())
	rec := db.Record{Key:key, URL:url, Expiration:expiration}
	err := s.dataSource.Save(rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *Service) randKey() string {
	n := config.Settings.KeyLength
	b := make([]byte, n)
	for i, cache, remain := n - 1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (s *Service) ClearRecordsAsync() {
	minutes := config.Settings.ClearTimeMinutes
	waitTime := time.Duration(minutes) * time.Minute
	for {
		time.Sleep(waitTime)
		now := uint64(time.Now().Unix())
		removed, err := s.dataSource.DeleteAllAfter(now)
		if err != nil {
			log.Println("[ERROR]", err)
		}
		log.Println("[INFO] Expired records removed. Count:", removed)
	}
}

func (s *Service) FindByKey(key string) *db.Record {
	return s.dataSource.Find(key)
}