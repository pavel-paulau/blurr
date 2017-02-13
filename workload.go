package qb

import (
	"context"
	"sync"
	"time"
)

type insertFn func(workerID int64, key string, value interface{}) error

type queryFn func(workerID int64, field string, arg interface{}) error

func singleLoad(wg *sync.WaitGroup, workerID int64, w *WorkloadSettings) {
	defer wg.Done()

	for payload := range generatePayload(workerID, w) {
		if err := w.IFn(workerID, payload.key, payload.value); err != nil {
			logFatalln(err)
		}
	}
}

// WorkloadSettings incorporates all possible workload settings.
type WorkloadSettings struct {
	NumWorkers, NumDocs, DocSize int64
	InsertPercentage             int
	Time                         time.Duration
	IFn                          insertFn
	QFn                          queryFn
	Hostname                     string
	Consistency                  string
	QueryType                    int
}

//
func (w *WorkloadSettings) SetQueryType(workload string) {
	switch workload {
	case "Q1":
		w.QueryType = q1query
	case "Q2":
		w.QueryType = q2query
	case "Q3":
		w.QueryType = q3query
	}
}

// Load executes the load phase - insertion of brand new items.
func Load(w *WorkloadSettings) {
	wg := sync.WaitGroup{}

	w.NumDocs /= w.NumWorkers

	for i := int64(0); i < w.NumWorkers; i++ {
		wg.Add(1)
		go singleLoad(&wg, i, w)
	}

	wg.Wait()
}

func singleRun(wg *sync.WaitGroup, workerID int64, w *WorkloadSettings, ctx context.Context) {
	defer wg.Done()

	ch1, ch2 := generateMixedPayload(w)

	for {
		select {
		case payload := <-ch1:
			if err := w.IFn(workerID, payload.key, payload.value); err != nil {
				logFatalln(err)
			}
		case payload := <-ch2:
			if err := w.QFn(workerID, payload.field, payload.arg); err != nil {
				logFatalln(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

var mu sync.RWMutex

// Run executes mixed workloads - a mix of queries and insert operations.
func Run(w *WorkloadSettings) {
	ctx, cancel := context.WithTimeout(context.Background(), w.Time)
	defer cancel()

	wg := sync.WaitGroup{}

	mu = sync.RWMutex{}
	currDocuments = w.NumDocs + 1

	for i := int64(0); i < w.NumWorkers; i++ {
		wg.Add(1)
		go singleRun(&wg, i, w, ctx)
	}

	wg.Wait()
}
