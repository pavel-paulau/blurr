package qb

import (
	"reflect"
	"testing"
)

func TestQ1(t *testing.T) {
	payload := q1(123456789, newZipf(123456789, 0))

	var expectedQueryType = "Q1"
	if payload.QueryType != expectedQueryType {
		t.Errorf("expected: %v, got: %v", expectedQueryType, payload.QueryType)
	}

	expectedProjection := []string{"address"}
	if !reflect.DeepEqual(payload.Projection, expectedProjection) {
		t.Errorf("expected: %v, got: %v", expectedProjection, payload.Projection)
	}

	expectedSelection := []Filter{{"_id", nil, false}}
	if payload.Selection[0].Field != expectedSelection[0].Field {
		t.Errorf("expected: %v, got: %v", expectedSelection[0].Field, payload.Selection[0].Field)
	}
}

func TestQ2(t *testing.T) {
	payload := q2(123456789, newZipf(123456789, 0))

	var expectedQueryType = "Q2"
	if payload.QueryType != expectedQueryType {
		t.Errorf("expected: %v, got: %v", expectedQueryType, payload.QueryType)
	}

	expectedProjection := []string{"address"}
	if !reflect.DeepEqual(payload.Projection, expectedProjection) {
		t.Errorf("expected: %v, got: %v", expectedProjection, payload.Projection)
	}

	expectedSelection := []Filter{{"email", nil, false}}
	if payload.Selection[0].Field != expectedSelection[0].Field {
		t.Errorf("expected: %v, got: %v", expectedSelection[0].Field, payload.Selection[0].Field)
	}
}

func TestQ3(t *testing.T) {
	payload := q3(123456789, newZipf(123456789, 0))

	var expectedQueryType = "Q3"
	if payload.QueryType != expectedQueryType {
		t.Errorf("expected: %v, got: %v", expectedQueryType, payload.QueryType)
	}

	expectedProjection := []string{"address"}
	if !reflect.DeepEqual(payload.Projection, expectedProjection) {
		t.Errorf("expected: %v, got: %v", expectedProjection, payload.Projection)
	}

	expectedSelection := []Filter{{"localgroup", nil, false}}
	if payload.Selection[0].Field != expectedSelection[0].Field {
		t.Errorf("expected: %v, got: %v", expectedSelection[0].Field, payload.Selection[0].Field)
	}
}

func TestQ4(t *testing.T) {
	payload := q4(123456789, newZipf(123456789, 0))

	var expectedQueryType = "Q4"
	if payload.QueryType != expectedQueryType {
		t.Errorf("expected: %v, got: %v", expectedQueryType, payload.QueryType)
	}

	expectedProjection := []string{"firstname", "lastname"}
	if !reflect.DeepEqual(payload.Projection, expectedProjection) {
		t.Errorf("expected: %v, got: %v", expectedProjection, payload.Projection)
	}

	expectedSelection := []Filter{{"address.zip", nil, false}, {"address.street", nil, true}}
	if payload.Selection[0].Field != expectedSelection[0].Field {
		t.Errorf("expected: %v, got: %v", expectedSelection[0].Field, payload.Selection[0].Field)
	}
}

func TestQ5(t *testing.T) {
	payload := q5(123456789, newZipf(123456789, 0))

	var expectedQueryType = "Q5"
	if payload.QueryType != expectedQueryType {
		t.Errorf("expected: %v, got: %v", expectedQueryType, payload.QueryType)
	}

	expectedProjection := []string{}
	if !reflect.DeepEqual(payload.Projection, expectedProjection) {
		t.Errorf("expected: %v, got: %v", expectedProjection, payload.Projection)
	}

	expectedSelection := []Filter{{"_id", nil, false}}
	if payload.Selection[0].Field != expectedSelection[0].Field {
		t.Errorf("expected: %v, got: %v", expectedSelection[0].Field, payload.Selection[0].Field)
	}
}
