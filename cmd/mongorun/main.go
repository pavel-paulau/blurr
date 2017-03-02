package main

import (
	"flag"
	"log"
	"runtime/debug"
	"time"

	"github.com/pavel-paulau/qb"
	"github.com/pavel-paulau/qb/mongo"
)

const _GOGC = 300

func init() {
	debug.SetGCPercent(_GOGC)
}

func main() {
	w := qb.WorkloadSettings{IFn: mongo.Insert, QFn: mongo.Query}
	var workload string

	flag.Int64Var(&w.NumWorkers, "workers", 1, "number of workload threads")
	flag.Int64Var(&w.NumDocs, "docs", 1e3, "number of documents to insert")
	flag.Int64Var(&w.DocSize, "size", 512, "average size of the documents")

	flag.IntVar(&w.InsertPercentage, "inserts", 5, "The percentage of insert operations")

	flag.DurationVar(&w.Time, "time", time.Minute, "Benchmark duration")

	flag.StringVar(&workload, "workload", "Q1", "Workload type")
	flag.StringVar(&w.Hostname, "hostname", "127.0.0.1", "MongoDB hostname")
	flag.BoolVar(&w.SSL, "ssl", false, "use SSL/TLS")

	flag.Parse()

	w.SetQueryType(workload)

	err := mongo.InitDatabase(w.Hostname, w.NumWorkers, w.SSL)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	qb.Run(&w)
}
