package db

type DataSource interface {
	Find(key string) *Record
	FindShort(url string) *Record
	Save(newRecord Record) (error)
	ExistsKey(key string) (bool, error)
	DeleteAllAfter(time int64) (removed int, err error)
}

type Record struct {
	Key string
	URL string
	Expiration int64
}

