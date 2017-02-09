package qb

import (
	"errors"
	"log"
	"testing"
)

func insertMock(workerID int64, key string, value interface{}) error {
	return nil
}

func insertMockFailure(workerID int64, key string, value interface{}) error {
	return errors.New("test")
}

func TestLoad(t *testing.T) {
	Load(insertMock, 1, 1e3, 1024)
}

func TestLoadFailure(t *testing.T) {
	defer func() {
		logFatalln = log.Fatalln
	}()

	var called bool
	logFatalln = func(v ...interface{}) {
		called = true
	}

	Load(insertMockFailure, 1, 1, 512)
	if !called {
		t.Fatal("the failure wasn't caught")
	}
}
