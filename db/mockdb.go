package db

import "sync"

type mockDataSource struct {
	sync.RWMutex
	records map[string]*Record
}

func (ds *mockDataSource) SaveWithExpiration(record *Record, expireIn int64) error {
	ds.Lock()
	ds.records[record.Key] = record
	ds.Unlock()
	return nil
}

func (ds *mockDataSource) Exists(key string) (bool, error) {
	ds.RLock()
	_, exists := ds.records[key]
	ds.RUnlock()
	return exists, nil
}

func (ds *mockDataSource) Find(key string) (*Record, error) {
	ds.RLock()
	rec := ds.records[key]
	ds.RUnlock()
	return rec, nil
}

func (ds *mockDataSource) Shutdown() { // skip
}
