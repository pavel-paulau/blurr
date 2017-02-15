package qb

import (
	"math/rand"
)

type kvPayload struct {
	key   string
	value *Doc
}

// Filter specifies selection filter
type Filter struct {
	Field  string
	Arg    interface{}
	IsText bool
}

// QueryPayload carries arguments and metadata required for query execution
type QueryPayload struct {
	QueryType  string
	Selection  []Filter
	Projection []string // Projection specifies the fields to return
}

const prefix = "user-profile"

const (
	insert = iota
	q1query
	q2query
	q3query
	q4query
	q5query
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func generatePayload(workerID int64, w *WorkloadSettings) chan kvPayload {
	payload := make(chan kvPayload, min(1e3, w.NumDocs))

	go func() {
		defer close(payload)

		for i := int64(0); i < w.NumDocs; i++ {
			j := workerID*w.NumDocs + i
			key := newKey(prefix, j)
			doc := newDoc(j, key, w.DocSize)
			payload <- kvPayload{key, &doc}
		}
	}()

	return payload
}

func initOpSet(insertPercentage, queryType int) []int {
	var operations []int
	for i := 0; i < insertPercentage; i++ {
		operations = append(operations, insert)
	}
	for i := 0; i < (100 - insertPercentage); i++ {
		operations = append(operations, queryType)
	}
	return operations
}

func generateSeq(insertPercentage, queryType int) chan int {
	opSeq := make(chan int, 1e3)

	opSet := initOpSet(insertPercentage, queryType)

	go func() {
		defer close(opSeq)

		for {
			for _, i := range rand.Perm(len(opSet)) {
				opSeq <- opSet[i]
			}
		}
	}()

	return opSeq
}

var currDocuments int64

func generateMixedPayload(w *WorkloadSettings) (chan kvPayload, chan *QueryPayload) {
	var keySpace int64

	ch1 := make(chan kvPayload, 1e3)
	ch2 := make(chan *QueryPayload, 1e3)

	go func() {
		defer close(ch1)
		defer close(ch2)

		for op := range generateSeq(w.InsertPercentage, w.QueryType) {
			keySpace = currDocuments

			switch op {
			case insert:
				mu.Lock()
				currDocuments++
				keySpace = currDocuments
				mu.Unlock()

				key := newKey(prefix, keySpace)
				doc := newDoc(keySpace, key, w.DocSize)
				ch1 <- kvPayload{key, &doc}
			case q1query:
				ch2 <- q1(currDocuments)
			case q2query:
				ch2 <- q2(currDocuments)
			case q3query:
				ch2 <- q3(currDocuments)
			case q4query:
				ch2 <- q4(currDocuments)
			case q5query:
				ch2 <- q5(currDocuments)
			}
		}
	}()

	return ch1, ch2
}
