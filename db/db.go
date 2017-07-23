package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	"github.com/nestorsokil/goto-url/config"
	"crypto/tls"
	"net"
)


var session *mgo.Session
var database string

func SetupConnection() *mgo.Session {
	var s *mgo.Session
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
		s, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
	} else {
		s, err = mgo.Dial(config.Settings.MongoUrls[0])
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
	}
	session = s
	database = config.Settings.Database
	return session
}

type Record struct {
	Key string
	URL string
	Expiration uint64
}

func Find(key string) *Record {
	database := session.DB(database)
	var result Record
	err := database.C("records").Find(bson.M{"key" : key}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func FindShort(url string) *Record {
	database := session.DB(database)
	var result Record
	err := database.C("records").Find(bson.M{"url" : url}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func Save(newRecord Record) (error) {
	database := session.DB(database)
	err := database.C("records").Insert(newRecord)
	if err != nil {
		return err
	}
	return nil
}

func ExistsKey(key string) (bool, error) {
	database := session.DB(database)
	count, err := database.C("records").Find(bson.M{"key": key}).Count()
	if err != nil {
		 return false, err
	}
	return count > 0, nil
}

func DeleteAllAfter(time uint64) (removed int, err error) {
	database := session.DB(database)
	info, err := database.C("records").RemoveAll(
		bson.M{"expiration": bson.M{"$gt": time}})
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}