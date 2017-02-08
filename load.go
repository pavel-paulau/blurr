package qb

import (
	"log"
	"sync"
)

type insertFn func(key string, value interface{}) error

var logFatalln = log.Fatalln

func singleLoad(wg *sync.WaitGroup, fn insertFn, workerID, docsPerWorker, docSize int64) {
	defer wg.Done()

	for kv := range generatePayload(workerID, docsPerWorker, docSize) {
		err := fn(kv.key, kv.value)
		if err != nil {
			logFatalln(err)
		}
	}
}

func Load(fn insertFn, numWorkers, numDocs, docSize int64) {
	wg := sync.WaitGroup{}

	for i := int64(0); i < numWorkers; i++ {
		wg.Add(1)
		go singleLoad(&wg, fn, i, numDocs, docSize)
	}

	wg.Wait()
}
