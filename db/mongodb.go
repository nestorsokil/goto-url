package db

import (
	"crypto/tls"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"

	"github.com/nestorsokil/goto-url/util"
)

type MongoDataSource struct {
	session  *mgo.Session
	database string
}

func (mongo *MongoDataSource) query() *mgo.Database {
	return mongo.session.DB(mongo.database)
}

func (mongo *MongoDataSource) Find(key string) (*Record, error) {
	var result Record
	err := mongo.query().C("records").Find(bson.M{"key": key}).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (mongo *MongoDataSource) FindShort(url string) (*Record, error) {
	var result Record
	err := mongo.query().C("records").Find(bson.M{"url": url}).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (mongo *MongoDataSource) Save(newRecord *Record) error {
	err := mongo.query().C("records").Insert(newRecord)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *MongoDataSource) ExistsKey(key string) (bool, error) {
	count, err := mongo.query().C("records").Find(bson.M{"key": key}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mongo *MongoDataSource) DeleteAllExpiredBefore(time int64) (removed int, err error) {
	info, err := mongo.query().C("records").RemoveAll(
		bson.M{"expiration": bson.M{"$lt": time}})
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}

func (mongo *MongoDataSource) Update(record *Record) error {
	return mongo.query().C("records").Update(bson.M{"key": record.Key}, record)
}

func (mongo *MongoDataSource) Shutdown() {
	mongo.session.Close()
}

func NewMongoDS(config *util.MongoConfig) (DataSource, error) {
	session, err := newMongoSession(config)
	if err != nil {
		return nil, err
	}
	return &MongoDataSource{session, config.DatabaseName}, nil
}

func newMongoSession(config *util.MongoConfig) (*mgo.Session, error) {
	var session *mgo.Session
	var err error
	if config.EnableTLS {
		tlsConfig := &tls.Config{}
		dialInfo := &mgo.DialInfo{
			Addrs:    config.MongoUrls,
			Username: config.MongoUser,
			Password: config.MongoPassword,
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return nil, err
		}
	} else {
		session, err = mgo.Dial(config.MongoUrls[0])
		if err != nil {
			return nil, err
		}
	}
	return session, nil
}
