package db

type MockDataSource struct {
	records map[string] *Record
}

func (ds *MockDataSource) Find(key string) *Record {
	return ds.records[key]
}

func (ds *MockDataSource) FindShort(url string) *Record {
	for _, rec := range ds.records {
		if rec.URL == url {
			return rec
		}
	}
	return nil
}

func (ds *MockDataSource) Save(newRecord Record) (error) {
	ds.records[newRecord.Key] = &newRecord
	return nil
}

func (ds *MockDataSource) ExistsKey(key string) (bool, error) {
	_, exists := ds.records[key]
	return exists, nil
}

func (ds *MockDataSource) DeleteAllAfter(time uint64) (removed int, err error) {
	count := 0
	for key, record := range ds.records {
		if record.Expiration < time {
			delete(ds.records, key)
			count++
		}
	}
	return count, nil
}

func NewMockDS() *MockDataSource {
	return &MockDataSource{make(map[string]*Record)}
}