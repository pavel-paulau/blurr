package qb

import (
	"testing"
)

func TestQ2(t *testing.T) {
	payload := q2(123456789)

	if payload.field != "email" {
		t.Errorf("expected: 'email', got: %v", payload.field)
	}
}

func TestQ3(t *testing.T) {
	payload := q3(123456789)

	if payload.field != "localgroup" {
		t.Errorf("expected: 'localgroup', got: %v", payload.field)
	}
}
