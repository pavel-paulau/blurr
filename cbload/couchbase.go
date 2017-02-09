package main

import (
	"fmt"

	"github.com/couchbase/go-couchbase"
)

var cbb *couchbase.Bucket

func initDatabase() error {
	baseURL := fmt.Sprintf("http://%s:8091/", hostname)

	c, err := couchbase.ConnectWithAuthCreds(baseURL, "default", "")
	if err != nil {
		return err
	}

	pool, err := c.GetPool("default")
	if err != nil {
		return err
	}

	cbb, err = pool.GetBucket("default")
	return err
}

func insert(workerID int64, key string, value interface{}) error {
	return cbb.Set(key, 0, value)
}
