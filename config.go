package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type address struct {
	Data string `yaml:"data"`
	N1QL string `yaml:"n1ql"`
}

type clientConfig struct {
	Address        address `yaml:"address"`
	Bucket         string  `yaml:"bucket"`
	BucketPassword string  `yaml:"bucket_password"`
}

type workloadConfig struct {
	CreatePercentage int   `yaml:"create_percentage"`
	ReadPercentage   int   `yaml:"read_percentage"`
	UpdatePercentage int   `yaml:"update_percentage"`
	DeletePercentage int   `yaml:"delete_percentage"`
	InitialDocuments int64 `yaml:"initial_documents"`
	Operations       int64 `yaml:"operations"`
	DocumentSize     int   `yaml:"document_size"`
	Workers          int64 `yaml:"workers"`
	Throughput       int64 `yaml:"throughput"`
	QueryWorkers     int64 `yaml:"query_workers"`
	QueryThroughput  int64 `yaml:"query_throughput"`
}

type queryConfig struct {
	Index       string `yaml:"index"`
	Consistency string `yaml:"consistency"`
}

type nbConfig struct {
	Database clientConfig   `yaml:"database"`
	Workload workloadConfig `yaml:"workload"`
	Query    queryConfig    `yaml:"query"`
}

func readConfig(path string) nbConfig {
	workload, err := ioutil.ReadFile(path)
	if err != nil {
		fatalf("error reading the configuration file '%v': %v\n", path, err)
	}

	var config nbConfig
	err = yaml.Unmarshal(workload, &config)
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
