package workloads

import (
	"fmt"
	"sync"
	"time"
)

type State struct {
	Operations, Records int64
	Errors              map[string]int
}

func (state *State) Init() {
	state.Errors = map[string]int{}
}

func (state *State) ReportThroughput(config Config, wg *sync.WaitGroup) {
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
