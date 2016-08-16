package main

import (
	"log"
	"sync"
	"time"
)

var (
	client   *Client
	config   Config
	state    *State
	workload *Workload
)

func init() {
	config = ReadConfig()

	client = newClient(config.Database)
	workload = &Workload{Config: config.Workload}

	state = newState(config.Workload.InitialDocuments)
}

func main() {
	wg := sync.WaitGroup{}
	wgStats := sync.WaitGroup{}

	for worker := 0; worker < config.Workload.Workers; worker++ {
		wg.Add(1)
		go workload.RunCRUDWorkload(client, state, &wg)
	}

	wgStats.Add(1)
	go state.ReportThroughput(config.Workload, &wgStats)

	if config.Workload.RunTime > 0 {
		time.Sleep(time.Duration(config.Workload.RunTime) * time.Second)
		log.Println("Shutting down workers")
	} else {
		wg.Wait()
		wgStats.Wait()
	}

	client.Shutdown()
}
