package db

import "sync"

type MockDataSource struct {
	sync.RWMutex
	records map[string]*Record
}

func (ds *MockDataSource) Find(key string) (*Record, error) {
	ds.RLock()
	rec := ds.records[key]
	ds.RUnlock()
	return rec, nil
}

func (ds *MockDataSource) FindShort(url string) (*Record, error) {
	ds.RLock()
	defer ds.RUnlock()
	for _, rec := range ds.records {
		if rec.URL == url {
			return rec, nil
		}
	}
	return nil, nil
}

func (ds *MockDataSource) Save(newRecord *Record) error {
	ds.Lock()
	ds.records[newRecord.Key] = newRecord
	ds.Unlock()
	return nil
}

func (ds *MockDataSource) ExistsKey(key string) (bool, error) {
	ds.RLock()
	_, exists := ds.records[key]
	ds.RUnlock()
	return exists, nil
}

func (ds *MockDataSource) DeleteAllExpiredBefore(time int64) (removed int, err error) {
	count := 0
	ds.RLock()
	defer ds.RUnlock()
	for key, record := range ds.records {
		if record.Expiration < time {
			ds.Lock()
			delete(ds.records, key)
			ds.Unlock()
			count++
		}
	}
	return count, nil
}

func (ds *MockDataSource) Update(record *Record) error {
	return ds.Save(record)
}

func (ds *MockDataSource) Shutdown() {
	// skip
}

func NewMockDS() (DataSource, error) {
	return &MockDataSource{records:make(map[string]*Record)}, nil
}
