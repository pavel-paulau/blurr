package main

import (
	"flag"
	"log"
	"runtime/debug"

	"github.com/pavel-paulau/qb"
	"github.com/pavel-paulau/qb/elastic"
)

const _GOGC = 300

func init() {
	debug.SetGCPercent(_GOGC)
}

func main() {
	w := qb.WorkloadSettings{IFn: elastic.Insert}

	flag.Int64Var(&w.NumWorkers, "workers", 1, "number of workload threads")
	flag.Int64Var(&w.NumDocs, "docs", 1e3, "number of documents to insert")
	flag.Int64Var(&w.DocSize, "size", 512, "average size of the documents")

	flag.StringVar(&w.Hostname, "hostname", "127.0.0.1", "Elasticsearch hostname")

	flag.Parse()

	err := elastic.InitDatabase(w.Hostname)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	qb.Load(&w)
}
