package db

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"crypto/tls"
	"net"
	"log"

	"github.com/nestorsokil/goto-url/config"
)

type MongoDataSource struct {
	session *mgo.Session
	database string
}

func(mongo *MongoDataSource) query() *mgo.Database {
	return mongo.session.DB(mongo.database)
}

func (mongo *MongoDataSource) Find(key string) *Record {
	var result Record
	err := mongo.query().C("records").Find(bson.M{"key" : key}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (mongo *MongoDataSource) FindShort(url string) *Record {
	var result Record
	err := mongo.query().C("records").Find(bson.M{"url" : url}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (mongo *MongoDataSource) Save(newRecord Record) (error) {
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

func (mongo *MongoDataSource) DeleteAllAfter(time uint64) (removed int, err error) {
	info, err := mongo.query().C("records").RemoveAll(
		bson.M{"expiration": bson.M{"$gt": time}})
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}

func NewMongoSession() *mgo.Session {
	var session *mgo.Session
	var err error
	if config.Settings.EnableTLS {
		tlsConfig := &tls.Config{}
		dialInfo := &mgo.DialInfo{
			Addrs: config.Settings.MongoUrls,
			Username: config.Settings.MongoUser,
			Password: config.Settings.MongoPassword,
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
	} else {
		session, err = mgo.Dial(config.Settings.MongoUrls[0])
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
	}
	return session
}

func NewMongoDS(session *mgo.Session, database string) *MongoDataSource {
	return &MongoDataSource{session, database}
}