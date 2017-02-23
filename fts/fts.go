package fts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/couchbase/go-couchbase"

	"github.com/pavel-paulau/qb"
)

const (
	bucketName = "default"
)

var (
	cbb      *couchbase.Bucket
	fts      *http.Client
	queryURL string
)

// InitDatabase initializes Couchbase Server client and HTTP client for FTS
// queries.
func InitDatabase(hostname string) error {
	baseURL := fmt.Sprintf("http://%s:8091/", hostname)

	c, err := couchbase.ConnectWithAuthCreds(baseURL, bucketName, "")
	if err != nil {
		return err
	}

	pool, err := c.GetPool("default")
	if err != nil {
		return err
	}

	cbb, err = pool.GetBucket(bucketName)
	if err != nil {
		return err
	}

	t := &http.Transport{MaxIdleConnsPerHost: 10240}
	fts = &http.Client{Transport: t}

	queryURL = fmt.Sprintf("http://%s:8094/api/index/default/query", hostname)

	return nil
}

// Insert adds new documents to Couchbase Server bucket using SET operation.
func Insert(_ int64, key string, value *qb.Doc) error {
	return cbb.Set(key, 0, value)
}

type ftsQuery struct {
	Fields []string          `json:"fields"`
	Query  map[string]string `json:"query"`
}

func executeQuery(q *ftsQuery) error {
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	j := bytes.NewReader(b)

	resp, err := fts.Post(queryURL, "application/json", j)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s - %d", "bad response", resp.StatusCode)
	}

	return nil
}

// Query finds matching documents using FTS API.
func Query(_ int64, payload *qb.QueryPayload) error {
	var query string
	for _, filter := range payload.Selection {
		query += " +" + filter.Field + ":" + filter.Arg.(string)
	}

	q := ftsQuery{
		Fields: payload.Projection,
		Query: map[string]string{
			"query": query,
		},
	}

	return executeQuery(&q)
}
