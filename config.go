package main

import (
	"encoding/json"
	"io/ioutil"
)

type clientConfig struct {
	Address        string `json:"address"`
	Bucket         string `json:"bucket"`
	BucketPassword string `json:"bucket_password"`
}

type workloadConfig struct {
	CreatePercentage int   `json:"create_percentage"`
	ReadPercentage   int   `json:"read_percentage"`
	UpdatePercentage int   `json:"update_percentage"`
	DeletePercentage int   `json:"delete_percentage"`
	InitialDocuments int64 `json:"initial_documents"`
	Operations       int64 `json:"operations"`
	DocumentSize     int   `json:"document_size"`
	Workers          int   `json:"workers"`
}

type nbConfig struct {
	Database clientConfig
	Workload workloadConfig
}

func readConfig(path string) nbConfig {
	workload, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config nbConfig
	err = json.Unmarshal(workload, &config)
	if err != nil {
		panic(err)
	}

	if config.Workload.ReadPercentage+config.Workload.UpdatePercentage+
		config.Workload.DeletePercentage > 0 && config.Workload.InitialDocuments == 0 {
		panic("Please specify non-zero 'initial_documents'")
	}

	if config.Workload.DocumentSize < sizeOverhead {
		panic("Document size must be greater than 450 bytes")
	}

	return config
}
