package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
)


var session *mgo.Session

func setupMongoSession() {
	s, err := mgo.Dial(conf.MongoUrl)
	if err != nil {
		log.Fatal("[FATAL] Cannot obtain session.")
	}
	session = s
}

type record struct {
	Key string
	URL string
}

func find(key string) *record {
	database := session.DB(conf.Database)
	var result record
	err := database.C("records").Find(bson.M{"key" : key}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func findShort(url string) *record {
	database := session.DB(conf.Database)
	var result record
	err := database.C("records").Find(bson.M{"url" : url}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func save(newRecord record) (error) {
	database := session.DB(conf.Database)
	err := database.C("records").Insert(newRecord)
	if err != nil {
		return err
	}
	return nil
}

func existsKey(key string) (bool, error) {
	database := session.DB(conf.Database)
	count, err := database.C("records").Find(bson.M{"key": key}).Count()
	if err != nil {
		 return false, err
	}
	return count > 0, nil
}