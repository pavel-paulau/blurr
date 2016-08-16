package main

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

type Client struct {
	Bucket *couchbase.Bucket
}

func newClient(config ClientConfig) *Client {
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

func (c *Client) Shutdown() {
	c.Bucket.Close()
}

func (c *Client) Create(key string, value map[string]interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *Client) Read(key string) error {
	result := map[string]interface{}{}
	return c.Bucket.Get(key, &result)
}

func (c *Client) Update(key string, value map[string]interface{}) error {
	return c.Bucket.Set(key, 0, value)
}

func (c *Client) Delete(key string) error {
	return c.Bucket.Delete(key)
}
