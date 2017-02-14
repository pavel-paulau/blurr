package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/pavel-paulau/qb"
)

var mgoc []*mgo.Session

const (
	dbName = "default"
	cName  = "default"
)

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
func Insert(workerId int64, key string, value *qb.Doc) error {
	value.ID = key
	session := mgoc[workerId]
	return session.DB("default").C("default").Insert(value)
}

// Query finds matching documents using MongoDB queries.
func Query(workerId int64, field string, arg interface{}) error {
	// FIXME: support multiple selectors
	query := bson.M{field: arg}

	// FIXME: support different projections
	projection := bson.M{"address": 1, "_id": 0}

	session := mgoc[workerId]
	var rs []interface{}

	return session.DB(dbName).C(cName).Find(query).Select(projection).All(&rs)
}
