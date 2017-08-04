package service

import (
	"testing"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/util"
	"time"
	"math"
)

func TestUrlService_GetRecord(t *testing.T) {
	conf := &util.Configuration{KeyLength: 5, ExpirationTimeHours: 1, }
	subject := New(db.NewMockDS(), conf)

	testUrl := "http://url.com"
	record, _ := subject.GetRecord(testUrl)
	key, url := record.Key, record.URL

	if url != testUrl {
		t.Errorf("Expected: %s, actual: %s", url, testUrl)
	}

	record, _ = subject.GetRecord(testUrl)
	key2, url2 := record.Key, record.URL

	if key != key2 {
		t.Errorf("Expected: %s, actual: %s", key, key2)
	}
	if url2 != testUrl {
		t.Errorf("Expected: %s, actual: %s", url2, testUrl)
	}

}

func TestService_FindByKey(t *testing.T) {
	conf := &util.Configuration{KeyLength: 5, ExpirationTimeHours: 1, }
	subject := New(db.NewMockDS(), conf)

	testUrl := "http://url.com"
	record, err := subject.GetRecord(testUrl)
	if err != nil {
		t.Error(err)
	}

	record = subject.FindByKey(record.Key)
	if record == nil {
		t.Errorf("Record was not found")
	}

	if record.URL != testUrl {
		t.Errorf("Expected: '%s', actual: '%s'", testUrl, record.URL)
	}

	now := time.Now().UnixNano()
	shouldExpireInApprox := now + conf.ExpirationTimeHours * time.Hour.Nanoseconds()
	countErr := 1 * time.Second.Nanoseconds()
	actualErr := int64(math.Abs(float64(record.Expiration - shouldExpireInApprox)))
	if actualErr > countErr {
		t.Errorf("Expected approximately: %d +- %d, actual %d",
			shouldExpireInApprox, countErr, record.Expiration)
	}
}

func TestUrlService_ClearRecordsAsync(t *testing.T) {
	conf := &util.Configuration{KeyLength: 5, ExpirationTimeHours: 0, ClearTimeSeconds: 1}
	subject := New(db.NewMockDS(), conf)

	record, _ := subject.GetRecord("http://url.com")
	shouldBeOk := subject.FindByKey(record.Key)

	if shouldBeOk == nil {
		t.Error("Record was not persisted.")
	}

	done := make(chan struct{})
	go subject.ClearRecordsAsync(done)
	time.Sleep(time.Duration((conf.ClearTimeSeconds+1) * time.Second.Nanoseconds()))
	done <- struct{}{}

	shouldBeNil := subject.FindByKey(record.Key)
	if shouldBeNil != nil {
		t.Errorf("Expected nil!")
	}
}


