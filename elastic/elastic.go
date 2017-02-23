package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pavel-paulau/qb"
)

var (
	elastic  *http.Client
	queryURL string
)

// InitDatabase initializes Couchbase Server client and HTTP client for FTS
// queries.
func InitDatabase(hostname string) error {
	t := &http.Transport{MaxIdleConnsPerHost: 10240}
	elastic = &http.Client{Transport: t}

	queryURL = fmt.Sprintf("http://%s:9200/default/default/", hostname)

	return nil
}

// Insert adds new documents to Couchbase Server bucket using SET operation.
func Insert(_ int64, key string, value *qb.Doc) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	j := bytes.NewReader(b)

	resp, err := elastic.Post(queryURL+key, "application/json", j)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode > 201 {
		return fmt.Errorf("%s - %d", "bad response", resp.StatusCode)
	}

	return nil
}

type filter struct {
	Term map[string]interface{} `json:"term"`
}

type boolFilter struct {
	Filter []filter `json:"filter"`
}

type query struct {
	Bool boolFilter `json:"bool"`
}

type elasticQuery struct {
	StoredFields []string `json:"stored_fields"`
	Query        query    `json:"query"`
}

func executeQuery(q *elasticQuery) error {
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	j := bytes.NewReader(b)

	resp, err := elastic.Post(queryURL+"_search", "application/json", j)
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

// Query finds matching documents using search API.
func Query(_ int64, payload *qb.QueryPayload) error {
	var filters []filter

	for _, selection := range payload.Selection {
		term := map[string]interface{}{
			selection.Field: selection.Arg,
		}
		filters = append(filters, filter{term})
	}

	q := elasticQuery{
		StoredFields: payload.Projection,
		Query: query{
			Bool: boolFilter{
				Filter: filters,
			},
		},
	}

	return executeQuery(&q)
}
