package main

import (
	"flag"
	"fmt"
	"sync"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: np workload.yaml")
	}
	flag.Parse()
	config := readConfig(flag.Arg(0))

	qClient := newQueryClient(config)
	dClient := newDataClient(config.Database)
	defer dClient.shutdown()

	workload := newWorkload(&config.Workload)
	payloads := workload.startPayloadFeed()
	go workload.reportThroughput()

	wg := sync.WaitGroup{}
	for worker := int64(0); worker < config.Workload.Workers; worker++ {
		wg.Add(1)
		go workload.runWorkload(dClient, payloads, &wg)
	}
	for worker := int64(0); worker < config.Workload.QueryWorkers; worker++ {
		go workload.runQueries(qClient)
	}
	wg.Wait()

	workload.reportLatency()
}
