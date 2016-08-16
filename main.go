package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	config := readConfig()

	client := newClient(config.Database)
	defer client.shutdown()

	workload := &Workload{Config: config.Workload}
	state := newState(config.Workload.InitialDocuments)

	wg := sync.WaitGroup{}
	wgStats := sync.WaitGroup{}

	for worker := 0; worker < config.Workload.Workers; worker++ {
		wg.Add(1)
		go workload.runCRUDWorkload(client, state, &wg)
	}

	wgStats.Add(1)
	go state.reportThroughput(config.Workload, &wgStats)

	if config.Workload.RunTime > 0 {
		time.Sleep(time.Duration(config.Workload.RunTime) * time.Second)
		log.Println("Shutting down workers")
	} else {
		wg.Wait()
		wgStats.Wait()
	}
}
