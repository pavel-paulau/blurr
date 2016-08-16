package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type clientConfig struct {
	Address        string
	Bucket         string
	BucketPassword string
}

type workloadConfig struct {
	CreatePercentage int
	ReadPercentage   int
	UpdatePercentage int
	DeletePercentage int
	InitialDocuments int64
	Operations       int64
	DocumentSize     int
	Workers          int
}

type Config struct {
	Database clientConfig
	Workload workloadConfig
}

func readConfig() (config Config) {
	flag.Usage = func() {
		fmt.Println("Usage: np workload.json")
	}
	flag.Parse()
	workloadPath := flag.Arg(0)

	workload, err := ioutil.ReadFile(workloadPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(workload, &config)
	if err != nil {
		panic(err)
	}

	if config.Workload.ReadPercentage+config.Workload.UpdatePercentage+
		config.Workload.DeletePercentage > 0 && config.Workload.InitialDocuments == 0 {
		panic("Please specify non-zero 'InitialDocuments'")
	}

	return
}
