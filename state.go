package main

import (
	"fmt"
	"time"
)

func reportThroughput(workload *Workload) {
	opsDone := int64(0)

	fmt.Println("Benchmark started:")
	for {
		time.Sleep(10 * time.Second)

		throughput := (workload.currentOperations - opsDone) / 10
		opsDone = workload.currentOperations

		fmt.Printf("%10v ops/sec; total operations: %v\n", throughput, opsDone)
	}
}
