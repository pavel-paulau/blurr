package main

import (
	"gopkg.in/mgo.v2"
)

var mgoc []*mgo.Session

func initDatabase() error {
	session, err := mgo.Dial(hostname)
	if err != nil {
		return err
	}

	for i := int64(0); i < numWorkers; i++ {
		mgoc = append(mgoc, session.New())
	}

	return nil
}

func insert(workerId int64, key string, value interface{}) error {
	session := mgoc[workerId]
	return session.DB("default").C("default").Insert(value)
}
