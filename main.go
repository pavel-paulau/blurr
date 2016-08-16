package main

import (
	"sync"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	config := readConfig()

	client := newClient(config.Database)
	defer client.shutdown()

	workload := newWorkload(&config.Workload)

	var opsBuffer int64 = min(1e7, config.Workload.Operations)
	ops := make(chan string, opsBuffer)
	go generateSeq(&config.Workload, ops)

	var payloadsBuffer int64 = min(1e6, opsBuffer)
	payloads := make(chan payload, payloadsBuffer)
	go workload.generatePayload(payloads, ops)

	go reportThroughput(workload)

	wg := sync.WaitGroup{}
	for worker := 0; worker < config.Workload.Workers; worker++ {
		wg.Add(1)
		go workload.runWorkload(client, payloads, &wg)
	}
	wg.Wait()
}
