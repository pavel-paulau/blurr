package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	statements map[string]string = map[string]string{
		"by_email": "SELECT * FROM `bucket-1` WHERE email='%s'",
	}
)

type queryClient struct {
	httpClient              *http.Client
	url, index, consistency string
}

func newQueryClient(config nbConfig) *queryClient {
	t := &http.Transport{MaxIdleConnsPerHost: int(config.Workload.QueryWorkers)}

	client := queryClient{
		httpClient:  &http.Client{Transport: t},
		url:         fmt.Sprintf("%s/query/service", config.Database.Address.N1QL),
		index:       config.Query.Index,
		consistency: config.Query.Consistency,
	}
	return &client
}

func (c *queryClient) post(statement string) error {
	values := url.Values{
		"statement":        []string{statement},
		"scan_consistency": []string{c.consistency},
	}

	resp, err := c.httpClient.PostForm(c.url, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	return nil
}

func (c *queryClient) query(value doc) error {
	var statement string

	switch c.index {
	case "by_email":
		statement = fmt.Sprintf(statements[c.index], value.Email)
	default:
		fatalf("unknown index: %s\n", c.index)
	}

	return c.post(statement)
}
