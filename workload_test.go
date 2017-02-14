package qb

import (
	"errors"
	"log"
	"testing"
	"time"
)

func insertMock(_ int64, _ string, _ *Doc) error {
	return nil
}

func insertMockFailure(_ int64, _ string, _ *Doc) error {
	return errors.New("test")
}

func queryMock(_ int64, _ string, _ interface{}) error {
	return nil
}

func TestLoad(t *testing.T) {
	w := WorkloadSettings{
		IFn:        insertMock,
		NumWorkers: 1,
		NumDocs:    1e3,
		DocSize:    1024,
	}
	Load(&w)
}

func TestLoadFailure(t *testing.T) {
	defer func() {
		logFatalln = log.Fatalln
	}()

	var called bool
	logFatalln = func(v ...interface{}) {
		called = true
	}

	w := WorkloadSettings{
		IFn:        insertMockFailure,
		NumWorkers: 1,
		NumDocs:    1,
		DocSize:    1024,
	}
	Load(&w)

	if !called {
		t.Fatal("the failure wasn't caught")
	}
}

func TestRun(t *testing.T) {
	w := WorkloadSettings{
		IFn:              insertMock,
		QFn:              queryMock,
		NumWorkers:       1,
		NumDocs:          1e3,
		Time:             10 * time.Millisecond,
		InsertPercentage: 50,
		DocSize:          512,
		QueryType:        q2query,
	}

	Run(&w)
}

func TestSetQuery(t *testing.T) {
	w := WorkloadSettings{}
	mapping := map[string]int{
		"Q1": q1query,
		"Q2": q2query,
		"Q3": q3query,
	}

	for k, v := range mapping {
		w.SetQueryType(k)
		if w.QueryType != v {
			t.Errorf("expected: %v, got: %v", v, w.QueryType)
		}
	}
}
