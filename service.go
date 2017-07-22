package main


func createNewRecord(url string) (*record, error) {
	var key string
	exists := true
	for exists {
		key = randKey(5)
		e, err := existsKey(key)
		if err != nil {
			return nil, err
		}
		exists = e
	}
	rec := record{randKey(5), url}
	err := save(rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}
