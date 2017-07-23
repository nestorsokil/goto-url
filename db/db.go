package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	"github.com/nestorsokil/goto-url/config"
)


var session *mgo.Session
var database string

func SetupConnection() *mgo.Session {
	s, err := mgo.Dial(config.Settings.MongoUrl)
	if err != nil {
		log.Fatal("[FATAL] Cannot obtain session.")
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