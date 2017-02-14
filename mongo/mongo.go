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
func Query(workerId int64, payload *qb.QueryPayload) error {
	query := bson.M{}
	for _, filter := range payload.Selection {
		if filter.IsText {
			query["$text"] = bson.M{"$search": filter.Arg}
		} else {
			query[filter.Field] = filter.Arg
		}
	}

	projection := bson.M{"_id": 0}
	for _, p := range payload.Projection {
		projection[p] = 1
	}

	session := mgoc[workerId]
	var rs []interface{}

	return session.DB(dbName).C(cName).Find(query).Select(projection).All(&rs)
}
