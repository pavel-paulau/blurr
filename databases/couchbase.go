package databases

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

type Couchbase struct {
	Bucket *couchbase.Bucket
}

func (cb *Couchbase) Init(config Config) {
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

	cb.Bucket = bucket
}

func (cb *Couchbase) Shutdown() {
	cb.Bucket.Close()
}

func (cb *Couchbase) Create(key string, value map[string]interface{}) error {
	return cb.Bucket.Set(key, 0, value)
}

func (cb *Couchbase) Read(key string) error {
	result := map[string]interface{}{}
	return cb.Bucket.Get(key, &result)
}

func (cb *Couchbase) Update(key string, value map[string]interface{}) error {
	return cb.Bucket.Set(key, 0, value)
}

func (cb *Couchbase) Delete(key string) error {
	return cb.Bucket.Delete(key)
}
