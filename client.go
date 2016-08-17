package main

import (
	"github.com/couchbaselabs/go-couchbase"
)

type dataClient struct {
	Bucket *couchbase.Bucket
}

func newDataClient(config clientConfig) *dataClient {
	conn, err := couchbase.Connect(config.Address.Data)
	if err != nil {
		fatalf("error connecting: %v\n", err)
	}

	pool, err := conn.GetPool("default")
	if err != nil {
		fatalf("error getting a pool: %v\n", err)
	}

	bucket, err := pool.GetBucketWithAuth(config.Bucket, config.Bucket, config.BucketPassword)
	if err != nil {
		fatalf("error getting a bucket: %v\n", err)
	}

	return &dataClient{bucket}
}

func (c *dataClient) shutdown() {
	c.Bucket.Close()
}

func (c *dataClient) create(key string, value interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *dataClient) read(key string) error {
	var result doc
	return c.Bucket.Get(key, &result)
}

func (c *dataClient) update(key string, value interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *dataClient) delete(key string) error {
	return c.Bucket.Delete(key)
}
