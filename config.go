package main

import (
	"encoding/json"
	"io/ioutil"
)

type address struct {
	Data string `json:"data"`
	N1QL string `json:"n1ql"`
}

type clientConfig struct {
	Address        address `json:"address"`
	Bucket         string  `json:"bucket"`
	BucketPassword string  `json:"bucket_password"`
}

type workloadConfig struct {
	CreatePercentage int   `json:"create_percentage"`
	ReadPercentage   int   `json:"read_percentage"`
	UpdatePercentage int   `json:"update_percentage"`
	DeletePercentage int   `json:"delete_percentage"`
	InitialDocuments int64 `json:"initial_documents"`
	Operations       int64 `json:"operations"`
	DocumentSize     int   `json:"document_size"`
	Workers          int64 `json:"workers"`
	Throughput       int64 `json:"throughput"`
}

type queryConfig struct {
	Index       string `json:"index"`
	Consistency string `json:"consistency"`
	Workers     int64  `json:"workers"`
}

type nbConfig struct {
	Database clientConfig   `json:"database"`
	Workload workloadConfig `json:"workload"`
	Query    queryConfig    `json:"query"`
}

func readConfig(path string) nbConfig {
	workload, err := ioutil.ReadFile(path)
	if err != nil {
		fatalf("error reading the configuration file '%v': %v\n", path, err)
	}

	var config nbConfig
	err = json.Unmarshal(workload, &config)
	if err != nil {
		fatalf("error parsing the configuration file %v: %v\n", path, err)
	}

	if config.Workload.ReadPercentage+config.Workload.UpdatePercentage+
		config.Workload.DeletePercentage > 0 && config.Workload.InitialDocuments == 0 {
		fatalln("please specify non-zero 'initial_documents'")
	}

	if config.Workload.DocumentSize < sizeOverhead {
		fatalln("document size must be greater than 450 bytes")
	}

	return config
}
