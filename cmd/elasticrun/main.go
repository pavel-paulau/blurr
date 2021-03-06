package main

import (
	"flag"
	"log"
	"runtime/debug"
	"time"

	"github.com/pavel-paulau/qb"
	"github.com/pavel-paulau/qb/elastic"
)

const _GOGC = 300

func init() {
	debug.SetGCPercent(_GOGC)
}

func main() {
	w := qb.WorkloadSettings{IFn: elastic.Insert, QFn: elastic.Query}
	var workload string

	flag.Int64Var(&w.NumWorkers, "workers", 1, "number of workload threads")
	flag.Int64Var(&w.NumDocs, "docs", 1e3, "number of documents to insert")
	flag.Int64Var(&w.DocSize, "size", 512, "average size of the documents")

	flag.IntVar(&w.InsertPercentage, "inserts", 5, "The percentage of insert operations")

	flag.DurationVar(&w.Time, "time", time.Minute, "Benchmark duration")

	flag.StringVar(&workload, "workload", "Q4", "Workload type")
	flag.StringVar(&w.Hostname, "hostname", "127.0.0.1", "Elasticsearch hostname")

	flag.Parse()

	w.SetQueryType(workload)

	err := elastic.InitDatabase(w.Hostname)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	qb.Run(&w)
}
