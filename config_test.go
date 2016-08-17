package main

import (
	"reflect"
	"testing"
)

func TestExampleConfig(t *testing.T) {
	actualConfig := readConfig("./example.json")

	expectedConfig := nbConfig{
		Database: clientConfig{
			Address: address{
				Data: "http://Administrator:password@127.0.0.1:8091",
				N1QL: "http://Administrator:password@127.0.0.1:8093",
			},
			Bucket:         "bucket-1",
			BucketPassword: "password",
		},
		Workload: workloadConfig{
			CreatePercentage: 2,
			ReadPercentage:   80,
			UpdatePercentage: 17,
			DeletePercentage: 1,
			InitialDocuments: 1000000,
			Operations:       1000000,
			DocumentSize:     1024,
			Workers:          100,
			Throughput:       100000,
			QueryWorkers:     10,
			QueryThroughput:  500,
		},
		Query: queryConfig{
			Index:       "by_email",
			Consistency: "request_plus",
		},
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("expected: %v, got: %v", expectedConfig, actualConfig)
	}
}
