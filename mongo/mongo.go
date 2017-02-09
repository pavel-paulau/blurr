package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoc []*mgo.Session

// InitDatabase initializes a pool of MongoDB clients.
func InitDatabase(hostname string, numWorkers int64) error {
	session, err := mgo.Dial(hostname)
	if err != nil {
		return err
	}

	for i := int64(0); i < numWorkers; i++ {
		mgoc = append(mgoc, session.New())
	}

	return nil
}

// Insert adds new documents to MongoDB collection.
func Insert(workerId int64, key string, value interface{}) error {
	session := mgoc[workerId]
	return session.DB("default").C("default").Insert(value)
}

// Query finds matching documents using MongoDB queries.
func Query(workerId int64, field string, arg interface{}) error {
	query := bson.M{field: arg}
	session := mgoc[workerId]
	var rs []interface{}
	// FIXME: support different queries
	return session.DB("default").C("default").Find(query).Select(bson.M{"address": 1}).All(&rs)
}
