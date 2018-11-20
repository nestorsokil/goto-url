package service

import (
	"github.com/nestorsokil/goto-url/conf"
	"github.com/nestorsokil/goto-url/db"
	"testing"
)

var c = &conf.TestConfig{Data: map[string]interface{}{
	conf.EnvStorage:          conf.InMemory,
	conf.EnvKeyLen:           10,
	conf.EnvExpirationMillis: 60000,
}}
var ds, _ = db.CreateStorage(c)
var subject = New(ds, c)

func TestUrlService_CreateRecordWithGeneratedKey(t *testing.T) {
	testUrl := "http://url.com"
	record, err := subject.CreateRecord(testUrl, "")
	if err != nil {
		t.Errorf("Error creating record %v", err)
	}
	if record.URL != testUrl {
		t.Errorf("Expected %v, Actual %v", testUrl, record.URL)
	}
	exists, err := ds.Exists(record.Key)
	if err != nil {
		t.Errorf("Error checking record existence %v", err)
	}
	if !exists {
		t.Errorf("Record was not created in the data storage")
	}
	saved, err := ds.Find(record.Key)
	if err != nil {
		t.Errorf("Error getting record %v", err)
	}
	if saved.URL != testUrl {
		t.Errorf("Expected %v, Actual %v", testUrl, saved.URL)
	}
	if saved.Key != record.Key {
		t.Errorf("Expected %v, Actual %v", record.Key, saved.Key)
	}
}

func TestUrlService_CreateRecordWithCustomKey(t *testing.T) {
	testUrl := "http://url.com"
	customKey := "one-two-three"
	record, err := subject.CreateRecord(testUrl, customKey)
	if err != nil {
		t.Errorf("Error creating record %v", err)
	}
	if record.URL != testUrl {
		t.Errorf("Expected %v, Actual %v", testUrl, record.URL)
	}
	if record.Key != customKey {
		t.Errorf("Expected %v, Actual %v", customKey, record.Key)
	}
	exists, err := ds.Exists(record.Key)
	if err != nil {
		t.Errorf("Error checking record existence %v", err)
	}
	if !exists {
		t.Errorf("Record was not created in the data storage")
	}
	saved, err := ds.Find(record.Key)
	if err != nil {
		t.Errorf("Error getting record %v", err)
	}
	if saved.URL != testUrl {
		t.Errorf("Expected %v, Actual %v", testUrl, saved.URL)
	}
	if saved.Key != record.Key {
		t.Errorf("Expected %v, Actual %v", record.Key, saved.Key)
	}
}

func TestUrlService_CreateRecordWithInvalidUrl(t *testing.T) {
	record, err := subject.CreateRecord("not-a-url", "")
	if err != ErrNotUrl {
		t.Errorf("Expected %v", ErrNotUrl)
	}
	if record != nil {
		t.Error("Expected nil")
	}
}

func TestUrlService_GetByKey(t *testing.T) {
	key, url := "a-b-c", "test.com"
	err := ds.SaveWithExpiration(&db.Record{Key: key, URL: url}, 10000)
	if err != nil {
		t.Errorf("Error setting up test data %v", err)
	}
	record, err := subject.FindByKey(key)
	if err != nil {
		t.Errorf("Error getting record %v", err)
	}
	if record.URL != url {
		t.Errorf("Expected %v, Actual %v", record.URL, url)
	}
	if record.Key != key {
		t.Errorf("Expected %v, Actual %v", record.Key, key)
	}
}

func TestGetNonExisting(t *testing.T) {
	key := "does_not_exist"
	record, err := subject.FindByKey(key)
	if err != ErrNotFound {
		t.Errorf("Excepted %v", ErrNotFound)
	}
	if record != nil {
		t.Errorf("Expected nil")
	}
}

func TestValidateUrl(t *testing.T) {
	validated, ok := validateUrl("golang.org")
	if !ok {
		t.Error("Expected to be valid")
	}
	if validated != "http://golang.org" {
		t.Errorf("Wrong output: %v", validated)
	}
	validated, ok = validateUrl("not-a-url")
	if ok {
		t.Error("Expected to be invalid")
	}
	validated, ok = validateUrl("https://stackoverflow.com/questions/11227809/why-is-it-faster-to-process-a-sorted-array-than-an-unsorted-array/11227902#11227902")
	if !ok {
		t.Error("Expected to be valid")
	}
	if validated != "https://stackoverflow.com/questions/11227809/why-is-it-faster-to-process-a-sorted-array-than-an-unsorted-array/11227902#11227902" {
		t.Errorf("Wrong output: %v", validated)
	}
}
