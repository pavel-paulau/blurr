package qb

import (
	"testing"
)

func TestGeneratePayloadSmall(t *testing.T) {
	w := WorkloadSettings{NumDocs: 100, DocSize: 512}

	for payload := range generatePayload(10, &w) {
		if payload.value.Age > 100 {
			t.Fatalf("unexpected payload: %+v", payload)
		}
	}
}

func TestGeneratePayloadLarge(t *testing.T) {
	w := WorkloadSettings{NumDocs: 1e4, DocSize: 256}

	for payload := range generatePayload(10, &w) {
		if payload.value.Notes != "" {
			t.Fatalf("unexpected payload: %+v", *payload.value)
		}
	}
}
