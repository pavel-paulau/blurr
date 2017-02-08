package qb

import (
	"testing"
)

func TestGeneratePayloadSmall(t *testing.T) {
	for payload := range generatePayload(10, 100, 512) {
		if payload.value.Age > 100 {
			t.Fatalf("unexpected payload: %+v", payload)
		}
	}
}

func TestGeneratePayloadLarge(t *testing.T) {
	for payload := range generatePayload(10, 1e4, 256) {
		if payload.value.Notes != "" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
	}
}
