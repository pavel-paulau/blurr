package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pavel-paulau/nb/databases"
	"github.com/pavel-paulau/nb/workloads"
)

type Config struct {
	Database databases.Config
	Workload workloads.Config
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
		config.Workload.DeletePercentage > 0 && config.Workload.Records == 0 {
		log.Fatal("Please specify non-zero 'Records'")
	}

	if config.Workload.Workers > 0 {
		config.Workload.Throughput /= config.Workload.Workers
	}

	return
}
