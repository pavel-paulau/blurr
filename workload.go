package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/samuel/go-metrics/metrics"
)

const (
	batchSize         = 100
	sizeOverhead  int = 450
	reservoirSize     = 1e5
)

var (
	defaultPercentiles     = []float64{0.5, 0.8, 0.9, 0.95, 0.99, 0.999}
	defaultPercentileNames = []string{"p50", "p80", "p90", "p95", "p99", "p999"}
)

type dbWorkload struct {
	config               *workloadConfig
	currentOperations    int64
	currentDocuments     int64
	deletedDocuments     int64
	queryLatency         metrics.Sample
	targetBatchTime      time.Duration
	targetQueryBatchTime time.Duration
}

func newWorkload(config *workloadConfig) *dbWorkload {
	w := dbWorkload{
		config:           config,
		currentDocuments: config.InitialDocuments,
		queryLatency:     metrics.NewUniformSample(reservoirSize),
	}

	if config.Throughput > 0 {
		throughput := config.Throughput / config.Workers
		w.targetBatchTime = batchSize * time.Duration(1e9/throughput)
	}

	if config.QueryThroughput > 0 {
		throughput := config.QueryThroughput / config.QueryWorkers
		w.targetQueryBatchTime = batchSize * time.Duration(1e9/throughput)
	}

	return &w
}

func (w *dbWorkload) startPayloadFeed() chan payload {
	var opsBuffer int64 = min(1e7, w.config.Operations)
	ops := make(chan string, opsBuffer)
	go generateSeq(w.config, ops)

	var payloadsBuffer int64 = min(1e6, opsBuffer)
	payloads := make(chan payload, payloadsBuffer)
	go w.generatePayload(payloads, ops)

	return payloads
}

func (w *dbWorkload) generateNewKey() string {
	w.currentDocuments++
	return fmt.Sprintf("%012d", w.currentDocuments)
}

func (w *dbWorkload) generateExistingKey() string {
	randRecord := 1 + rand.Int63n(w.currentDocuments-w.deletedDocuments)
	randRecord += w.deletedDocuments
	return fmt.Sprintf("%012d", randRecord)
}

func (w *dbWorkload) generateKeyForRemoval() string {
	w.deletedDocuments++
	return fmt.Sprintf("%012d", w.deletedDocuments)
}

func (w *dbWorkload) generateValue(key string) doc {
	return newDoc(key, w.config.DocumentSize)
}

func initOpsSet(config *workloadConfig) []string {
	operations := []string{}
	for i := 0; i < config.CreatePercentage; i++ {
		operations = append(operations, "c")
	}
	for i := 0; i < config.ReadPercentage; i++ {
		operations = append(operations, "r")
	}
	for i := 0; i < config.UpdatePercentage; i++ {
		operations = append(operations, "u")
	}
	for i := 0; i < config.DeletePercentage; i++ {
		operations = append(operations, "d")
	}
	if len(operations) != 100 {
		panic("wrong workload configuration: sum of percentages is not equal 100")
	}
	return operations
}

func generateSeq(config *workloadConfig, ops chan string) {
	defer close(ops)

	opsSet := initOpsSet(config)

	for {
		for _, i := range rand.Perm(len(opsSet)) {
			ops <- opsSet[i]
		}

		config.Operations -= int64(len(opsSet))
		if config.Operations == 0 {
			break
		}
	}
}

type payload struct {
	op, key string
	value   doc
}

func (w *dbWorkload) generatePayload(payloads chan payload, ops chan string) {
	defer close(payloads)

	for op := range ops {
		var key string
		var value doc

		switch op {
		case "c":
			key = w.generateNewKey()
			value = w.generateValue(key)
		case "r":
			key = w.generateExistingKey()
		case "u":
			key = w.generateExistingKey()
			value = w.generateValue(key)
		case "d":
			key = w.generateKeyForRemoval()
		}

		payloads <- payload{op, key, value}
	}
}

func (w *dbWorkload) doData(client *dataClient, p payload) {
	var err error

	switch p.op {
	case "c":
		err = client.create(p.key, p.value)
	case "r":
		err = client.read(p.key)
	case "u":
		err = client.update(p.key, p.value)
	case "d":
		err = client.delete(p.key)
	}

	if err != nil {
		log.Println(err)
	}
}

func (w *dbWorkload) sleep(t0 *time.Time, targetTime time.Duration) {
	batchTime := time.Now().Sub(*t0)
	sleepTime := targetTime - batchTime
	if sleepTime > 0 {
		time.Sleep(time.Duration(sleepTime))
	}
	*t0 = time.Now()
}

func (w *dbWorkload) runWorkload(client *dataClient, payloads chan payload, wg *sync.WaitGroup) {
	defer wg.Done()

	var batch int64
	t0 := time.Now()

	for p := range payloads {
		w.currentOperations++
		w.doData(client, p)

		if w.config.Throughput > 0 {
			batch++
			if batch == batchSize {
				w.sleep(&t0, w.targetBatchTime)
				batch = 0
			}
		}
	}
}

func (w *dbWorkload) doQuery(client *queryClient) {
	key := w.generateExistingKey()
	value := w.generateValue(key)

	t0 := time.Now()
	if err := client.query(value); err == nil {
		latency := time.Now().Sub(t0)
		w.queryLatency.Update(int64(latency))
	} else {
		log.Println(err)
	}
}

func (w *dbWorkload) runQueries(client *queryClient) {
	var batch int64
	t0 := time.Now()

	for {
		w.doQuery(client)

		if w.config.QueryThroughput > 0 {
			batch++
			if batch == batchSize {
				w.sleep(&t0, w.targetQueryBatchTime)
				batch = 0
			}
		}
	}
}

func (w *dbWorkload) reportThroughput() {
	var opsDone int64

	fmt.Println("Workload started.")
	for {
		time.Sleep(5 * time.Second)

		opsThroughput := (w.currentOperations - opsDone) / 5
		opsDone = w.currentOperations

		fmt.Printf("%10d ops/sec; total ops: %s;\n", opsThroughput, humanize.Comma(w.currentOperations))
	}
}

func (w *dbWorkload) reportLatency() {
	histogram := metrics.NewSampledHistogram(w.queryLatency)
	percentiles := histogram.Percentiles(defaultPercentiles)

	fmt.Println("Query latency:")
	for i, p := range defaultPercentileNames {
		fmt.Printf("\t%s\t: %s us\n", p, humanize.Comma(percentiles[i]/1e3))
	}
}
