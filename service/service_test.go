package service

import (
	"testing"
	"github.com/nestorsokil/goto-url/db"
)

var service = New(&db.MockDataSource{})

//TODO : for tests to work config has to be refactored
func TestService_FindByKey(t *testing.T) {
	testUrl := "http://url.com"
	record, _ := service.GetRecord(testUrl)
	record = service.FindByKey(record.Key)

	if record.URL != testUrl {
		t.Fail()
	}
}


