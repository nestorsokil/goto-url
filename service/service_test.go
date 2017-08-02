package service

import (
	"testing"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/config"
)

var subject = New(db.NewMockDS(),
	&config.Config{
		KeyLength:5,
	})

func TestService_FindByKey(t *testing.T) {
	testUrl := "http://url.com"
	record, _ := subject.GetRecord(testUrl)
	record = subject.FindByKey(record.Key)

	if record.URL != testUrl {
		t.Fail()
	}
}


