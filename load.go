package qb

import (
	"log"
	"sync"
)

type insertFn func(workerID int64, key string, value interface{}) error

var logFatalln = log.Fatalln

func singleLoad(wg *sync.WaitGroup, fn insertFn, workerID, numDocs, docSize int64) {
	defer wg.Done()

	for kv := range generatePayload(workerID, numDocs, docSize) {
		err := fn(workerID, kv.key, kv.value)
		if err != nil {
			logFatalln(err)
		}
	}
}

func Load(fn insertFn, numWorkers, numDocs, docSize int64) {
	wg := sync.WaitGroup{}

	docsPerWorker := numDocs / numWorkers

	for i := int64(0); i < numWorkers; i++ {
		wg.Add(1)
		go singleLoad(&wg, fn, i, docsPerWorker, docSize)
	}

	wg.Wait()
}
