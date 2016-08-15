package workloads

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/pavel-paulau/nb/databases"
)

type State struct {
	Operations, Records int64
	Errors              map[string]int
	Events              map[string]time.Time
	Latency             map[string][]float64
}

func (state *State) Init() {
	state.Errors = map[string]int{}
	state.Events = map[string]time.Time{}
	state.Latency = map[string][]float64{
		"Create": {},
		"Read":   {},
		"Update": {},
		"Delete": {},
	}
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

func (state *State) MeasureLatency(database databases.Database,
	workload Workload, config Config, wg *sync.WaitGroup) {
	defer wg.Done()

	for state.Operations < config.Operations {
		if config.CreatePercentage > 0 {
			state.Operations++
			state.Records++
			key := workload.GenerateNewKey(state.Records)
			value := workload.GenerateValue(key, config.ValueSize)
			t0 := time.Now()
			database.Create(key, value)
			t1 := time.Now()
			latency := float64(t1.Sub(t0)/time.Microsecond) / 1000
			state.Latency["Create"] = append(state.Latency["Create"], latency)
		}
		if config.ReadPercentage > 0 {
			state.Operations++
			key := workload.GenerateExistingKey(state.Records)
			t0 := time.Now()
			database.Read(key)
			t1 := time.Now()
			latency := float64(t1.Sub(t0)/time.Microsecond) / 1000
			state.Latency["Read"] = append(state.Latency["Read"], latency)
		}
		if config.UpdatePercentage > 0 {
			state.Operations++
			key := workload.GenerateExistingKey(state.Records)
			value := workload.GenerateValue(key, config.ValueSize)
			t0 := time.Now()
			database.Update(key, value)
			t1 := time.Now()
			latency := float64(t1.Sub(t0)/time.Microsecond) / 1000
			state.Latency["Update"] = append(state.Latency["Update"], latency)
		}
		if config.DeletePercentage > 0 {
			state.Operations++
			key := workload.GenerateKeyForRemoval()
			t0 := time.Now()
			database.Delete(key)
			t1 := time.Now()
			latency := float64(t1.Sub(t0)/time.Microsecond) / 1000
			state.Latency["Delete"] = append(state.Latency["Delete"], latency)
		}
		time.Sleep(time.Second)
	}
}

func calcPercentile(data []float64, p float64) float64 {
	sort.Float64s(data)

	k := float64(len(data)-1) * p
	f := math.Floor(k)
	c := math.Ceil(k)
	if f == c {
		return data[int(k)]
	} else {
		return data[int(f)]*(c-k) + data[int(c)]*(k-f)
	}
}

func calcMean(data []float64) float64 {
	sum := float64(0)
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func (state *State) ReportSummary() {
	for _, op := range []string{"Create", "Read", "Update", "Delete"} {
		if len(state.Latency[op]) > 0 {
			fmt.Printf("%v latency:\n", op)
			for _, percentile := range []float64{0.8, 0.9, 0.95, 0.99} {
				value := calcPercentile(state.Latency[op], percentile)
				fmt.Printf("\t%vth percentile: %.2f ms\n", percentile*100, value)
			}
			value := calcMean(state.Latency[op])
			fmt.Printf("\tMean: %.2f ms\n", value)
		}
	}

	if len(state.Errors) > 0 {
		fmt.Println("Errors:")
		fmt.Printf("\tCreate : %v\n", state.Errors["c"])
		fmt.Printf("\tRead   : %v\n", state.Errors["r"])
		fmt.Printf("\tUpdate : %v\n", state.Errors["u"])
		fmt.Printf("\tDelete : %v\n", state.Errors["d"])
		fmt.Printf("\tTotal  : %v\n", state.Errors["total"])
	}
	fmt.Printf("Time elapsed:\n\t%v\n",
		state.Events["Finished"].Sub(state.Events["Started"]))
}
