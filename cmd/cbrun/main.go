package main

import (
	"flag"
	"log"
	"runtime/debug"
	"time"

	"github.com/pavel-paulau/qb"
	"github.com/pavel-paulau/qb/couchbase"
)

const _GOGC = 300

func init() {
	debug.SetGCPercent(_GOGC)
}

func main() {
	w := qb.WorkloadSettings{IFn: cb.Insert, QFn: cb.Query}
	var workload string

	flag.Int64Var(&w.NumWorkers, "workers", 1, "number of workload threads")
	flag.Int64Var(&w.NumDocs, "docs", 1e3, "number of documents to insert")
	flag.Int64Var(&w.DocSize, "size", 512, "average size of the documents")

	flag.IntVar(&w.InsertPercentage, "inserts", 5, "The percentage of insert operations")

	flag.DurationVar(&w.Time, "time", time.Minute, "Benchmark duration")

	flag.StringVar(&workload, "workload", "Q1", "Workload type")
	flag.StringVar(&w.Hostname, "hostname", "127.0.0.1", "Couchbase Server hostname")
	flag.StringVar(&w.Consistency, "consistency", "not_bounded", "N1QL scan consistency")

	flag.Parse()

	w.SetQueryType(workload)

	err := cb.InitDatabase(w.Hostname, w.Consistency)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	qb.Run(&w)
}
