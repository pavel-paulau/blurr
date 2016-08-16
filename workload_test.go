package main

import (
	"testing"
)

func TestOpsChannel(t *testing.T) {
	var totalOps int64 = 1000

	config := workloadConfig{
		CreatePercentage: 10,
		ReadPercentage:   20,
		UpdatePercentage: 30,
		DeletePercentage: 40,
		Operations:       totalOps,
	}

	ops := make(chan string, totalOps)
	go generateSeq(&config, ops)

	counter := map[string]int64{}
	for op := range ops {
		counter[op]++
	}

	expectedCreates := totalOps * int64(config.CreatePercentage) / 100
	if counter["c"] != expectedCreates {
		t.Errorf("expected: %v creates, got: %v", expectedCreates, counter["c"])
	}

	expectedReads := totalOps * int64(config.ReadPercentage) / 100
	if counter["r"] != expectedReads {
		t.Errorf("expected: %v reads, got: %v", expectedReads, counter["r"])
	}

	expectedUpdates := totalOps * int64(config.UpdatePercentage) / 100
	if counter["u"] != expectedUpdates {
		t.Errorf("expected: %v updates, got: %v", expectedUpdates, counter["u"])
	}

	expectedDeletes := totalOps * int64(config.DeletePercentage) / 100
	if counter["d"] != expectedDeletes {
		t.Errorf("expected: %v deletes, got: %v", expectedDeletes, counter["d"])
	}
}

func TestPayloadChannel(t *testing.T) {
	var totalOps int64 = 1000

	config := workloadConfig{
		CreatePercentage: 10,
		ReadPercentage:   20,
		UpdatePercentage: 30,
		DeletePercentage: 40,
		InitialDocuments: 1000,
		Operations:       totalOps,
		DocumentSize:     1024,
	}

	workload := newWorkload(&config)

	ops := make(chan string, totalOps)
	go generateSeq(&config, ops)

	payloads := make(chan payload, totalOps)
	go workload.generatePayload(payloads, ops)

	var currOps int64
	for range payloads {
		currOps++
	}

	if currOps != totalOps {
		t.Fatalf("expected %v ops, got: %v", totalOps, currOps)
	}
}

func TestPanicOps(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("generator did not panic")
		}
	}()

	config := workloadConfig{
		CreatePercentage: 0,
		ReadPercentage:   0,
		UpdatePercentage: 0,
		DeletePercentage: 0,
		Operations:       1000,
	}

	ops := make(chan string, 1000)
	generateSeq(&config, ops)
}
