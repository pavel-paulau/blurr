package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/pavel-paulau/nb/databases"
	"github.com/pavel-paulau/nb/workloads"
)

var config Config
var database databases.Database
var workload workloads.Workload
var state workloads.State

func init() {
	config = ReadConfig()

	database = &databases.Couchbase{}

	r := rand.New(rand.NewSource(0))
	zipf := rand.NewZipf(r, 1.4, 9.0, 1000)
	workload = &workloads.N1QL{
		Config: config.Workload,
		Zipf:   *zipf,
	}
	workload.SetImplementation(workload)

	database.Init(config.Database)

	state = workloads.State{}
	state.Records = config.Workload.Records
	state.Init()
}

func main() {
	wg := sync.WaitGroup{}
	wgStats := sync.WaitGroup{}

	state.Events["Started"] = time.Now()
	for worker := 0; worker < config.Workload.Workers; worker++ {
		wg.Add(1)
		go workload.RunCRUDWorkload(database, &state, &wg)
	}

	wgStats.Add(2)
	go state.ReportThroughput(config.Workload, &wgStats)
	go state.MeasureLatency(database, workload, config.Workload, &wgStats)

	if config.Workload.RunTime > 0 {
		time.Sleep(time.Duration(config.Workload.RunTime) * time.Second)
		log.Println("Shutting down workers")
	} else {
		wg.Wait()
		wgStats.Wait()
	}

	database.Shutdown()
	state.Events["Finished"] = time.Now()
	state.ReportSummary()
}
