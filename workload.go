package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

const (
	sizeOverhead int = 450
)

type Workload struct {
	config            *workloadConfig
	currentOperations int64
	currentDocuments  int64
	deletedDocuments  int64
}

func newWorkload(config *workloadConfig) *Workload {
	return &Workload{
		config:           config,
		currentDocuments: config.InitialDocuments,
	}
}

func (w *Workload) generateNewKey() string {
	w.currentDocuments++
	return fmt.Sprintf("%012d", w.currentDocuments)
}

func (w *Workload) generateExistingKey() string {
	randRecord := 1 + rand.Int63n(w.currentDocuments-w.deletedDocuments)
	randRecord += w.deletedDocuments
	return fmt.Sprintf("%012d", randRecord)
}

func (w *Workload) generateKeyForRemoval() string {
	w.deletedDocuments++
	return fmt.Sprintf("%012d", w.deletedDocuments)
}

func (w *Workload) generateValue(key string) doc {
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

func (w *Workload) generatePayload(payloads chan payload, ops chan string) {
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

func (w *Workload) do(client *Client, p payload) {
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

func (w *Workload) runWorkload(client *Client, payloads chan payload, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range payloads {
		w.currentOperations++
		w.do(client, p)
	}
}
