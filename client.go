package main

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

type Client struct {
	Bucket *couchbase.Bucket
}

func newClient(config clientConfig) *Client {
	conn, err := couchbase.Connect(config.Address)
	if err != nil {
		log.Fatalf("error connecting: %v", err)
	}

	pool, err := conn.GetPool("default")
	if err != nil {
		log.Fatalf("error getting a pool: %v", err)
	}

	bucket, err := pool.GetBucketWithAuth(config.Bucket, config.Bucket, config.BucketPassword)
	if err != nil {
		log.Fatalf("error getting a bucket: %v", err)
	}

	return &Client{bucket}
}

func (c *Client) shutdown() {
	c.Bucket.Close()
}

func (c *Client) create(key string, value interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *Client) read(key string) error {
	var result doc
	return c.Bucket.Get(key, &result)
}

func (c *Client) update(key string, value interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *Client) delete(key string) error {
	return c.Bucket.Delete(key)
}
