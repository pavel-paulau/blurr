package main

import (
	"fmt"
	"sync"
	"time"
)

type State struct {
	Operations, Documents int64
	Errors                map[string]int
}

func newState(initialDocuments int64) *State {
	return &State{Errors: map[string]int{}, Documents: initialDocuments}
}

func (state *State) reportThroughput(config workloadConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	opsDone := int64(0)
	samples := 1
	fmt.Println("Benchmark started:")
	for state.Operations < config.Operations {
		time.Sleep(10 * time.Second)
		throughput := (state.Operations - opsDone) / 10
		opsDone = state.Operations
		fmt.Printf("%6v seconds: %10v ops/sec; total operations: %v; total errors: %v\n",
			samples*10, throughput, opsDone, state.Errors["total"])
		samples++
	}
}
