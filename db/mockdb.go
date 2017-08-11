package db

type MockDataSource struct {
	records map[string]*Record
}

func (ds *MockDataSource) Find(key string) (*Record, error) {
	return ds.records[key], nil
}

func (ds *MockDataSource) FindShort(url string) (*Record, error) {
	for _, rec := range ds.records {
		if rec.URL == url {
			return rec, nil
		}
	}
	return nil, nil
}

func (ds *MockDataSource) Save(newRecord *Record) error {
	ds.records[newRecord.Key] = newRecord
	return nil
}

func (ds *MockDataSource) ExistsKey(key string) (bool, error) {
	_, exists := ds.records[key]
	return exists, nil
}

func (ds *MockDataSource) DeleteAllExpiredBefore(time int64) (removed int, err error) {
	count := 0
	for key, record := range ds.records {
		if record.Expiration < time {
			delete(ds.records, key)
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
	return &MockDataSource{make(map[string]*Record)}, nil
}
