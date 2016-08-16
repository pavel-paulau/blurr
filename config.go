package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

type ClientConfig struct {
	Address        string
	Bucket         string
	BucketPassword string
}

type WorkloadConfig struct {
	Type             string
	CreatePercentage int
	ReadPercentage   int
	UpdatePercentage int
	DeletePercentage int
	InitialDocuments int64
	Operations       int64
	ValueSize        int
	Workers          int
	Throughput       int
	RunTime          int
}

type Config struct {
	Database ClientConfig
	Workload WorkloadConfig
}

func ReadConfig() (config Config) {
	flag.Usage = func() {
		fmt.Println("Usage: np workload.conf")
	}
	flag.Parse()
	workload_path := flag.Arg(0)

	workload, err := ioutil.ReadFile(workload_path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(workload, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Workload.ReadPercentage+config.Workload.UpdatePercentage+
		config.Workload.DeletePercentage > 0 && config.Workload.InitialDocuments == 0 {
		log.Fatal("Please specify non-zero 'InitialDocuments'")
	}

	if config.Workload.Workers > 0 {
		config.Workload.Throughput /= config.Workload.Workers
	}

	return
}
