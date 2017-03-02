package main

import (
	"flag"
	"log"
	"runtime/debug"

	"github.com/pavel-paulau/qb"
	"github.com/pavel-paulau/qb/couchbase"
)

const _GOGC = 300

func init() {
	debug.SetGCPercent(_GOGC)
}

func main() {
	w := qb.WorkloadSettings{IFn: cb.Insert}

	flag.Int64Var(&w.NumWorkers, "workers", 1, "number of workload threads")
	flag.Int64Var(&w.NumDocs, "docs", 1e3, "number of documents to insert")
	flag.Int64Var(&w.DocSize, "size", 512, "average size of the documents")

	flag.StringVar(&w.Hostname, "hostname", "127.0.0.1", "Couchbase Server hostname")

	flag.Parse()

	err := cb.InitDatabase(w.Hostname, w.Consistency)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	qb.Load(&w)
}
