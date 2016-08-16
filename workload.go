package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	batchSize    int = 100
	sizeOverhead int = 450
)

type Workload struct {
	Config       workloadConfig
	DeletedItems int64
}

func (w *Workload) generateNewKey(currentDocuments int64) string {
	return fmt.Sprintf("%012d", currentDocuments)
}

func (w *Workload) generateExistingKey(currentDocuments int64) string {
	randRecord := 1 + rand.Int63n(currentDocuments-w.DeletedItems)
	randRecord += w.DeletedItems
	return fmt.Sprintf("%012d", randRecord)
}

func (w *Workload) generateKeyForRemoval() string {
	w.DeletedItems++
	return fmt.Sprintf("%012d", w.DeletedItems)
}

func (w *Workload) generateValue(key string, size int) doc {
	if size < sizeOverhead {
		log.Fatalf("Wrong workload configuration: minimal value size is %v", sizeOverhead)
	}

	return newDoc(key, size)
}

func (w *Workload) prepareBatch() []string {
	operations := make([]string, 0, batchSize)
	for i := 0; i < w.Config.CreatePercentage; i++ {
		operations = append(operations, "c")
	}
	for i := 0; i < w.Config.ReadPercentage; i++ {
		operations = append(operations, "r")
	}
	for i := 0; i < w.Config.UpdatePercentage; i++ {
		operations = append(operations, "u")
	}
	for i := 0; i < w.Config.DeletePercentage; i++ {
		operations = append(operations, "d")
	}
	if len(operations) != batchSize {
		log.Fatal("Wrong workload configuration: sum of percentages is not equal 100")
	}
	return operations
}

func (w *Workload) prepareSeq(size int64) chan string {
	operations := w.prepareBatch()
	seq := make(chan string, batchSize)
	go func() {
		for i := int64(0); i < size; i += int64(batchSize) {
			for _, randI := range rand.Perm(batchSize) {
				seq <- operations[randI]
			}
		}
	}()
	return seq
}

func (w *Workload) doBatch(client *Client, state *State, seq chan string) {
	for i := 0; i < batchSize; i++ {
		op := <-seq
		if state.Operations < w.Config.Operations {
			var err error
			state.Operations++
			switch op {
			case "c":
				state.Documents++
				key := w.generateNewKey(state.Documents)
				value := w.generateValue(key, w.Config.ValueSize)
				err = client.create(key, value)
			case "r":
				key := w.generateExistingKey(state.Documents)
				err = client.read(key)
			case "u":
				key := w.generateExistingKey(state.Documents)
				value := w.generateValue(key, w.Config.ValueSize)
				err = client.update(key, value)
			case "d":
				key := w.generateKeyForRemoval()
				err = client.delete(key)
			}
			if err != nil {
				state.Errors[op]++
				state.Errors["total"]++
			}
		}
	}
}

func (w *Workload) runWorkload(client *Client, state *State, wg *sync.WaitGroup, targetBatchTime float64, seq chan string) {
	for state.Operations < w.Config.Operations {
		t0 := time.Now()
		w.doBatch(client, state, seq)
		t1 := time.Now()

		if !math.IsInf(targetBatchTime, 0) {
			targetBatchTime := time.Duration(targetBatchTime * math.Pow10(9))
			actualBatchTime := t1.Sub(t0)
			sleepTime := (targetBatchTime - actualBatchTime)
			if sleepTime > 0 {
				time.Sleep(time.Duration(sleepTime))
			}
		}
	}
}

func (w *Workload) runCRUDWorkload(client *Client, state *State, wg *sync.WaitGroup) {
	defer wg.Done()

	seq := w.prepareSeq(w.Config.Operations)

	targetBatchTime := float64(batchSize) / float64(w.Config.Throughput)

	w.runWorkload(client, state, wg, targetBatchTime, seq)
}
