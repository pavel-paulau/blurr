package qb

import (
	"reflect"
	"testing"
)

func TestDoc(t *testing.T) {
	expectedDoc := doc{
		Name:         "ecdb3e e921c9",
		Email:        "3d13c6@a2d1f3.com",
		Street:       "400f1d0a",
		City:         "90ac48",
		County:       "40efd6",
		Country:      "1811db",
		State:        "WY",
		FullState:    "Montana",
		Realm:        "15e3f5",
		Coins:        213.54,
		Category:     1,
		Achievements: []int16{0, 135, 92},
		GMTime:       []int16{1972, 3, 3, 0, 0, 0, 4, 63, 0},
		Year:         1989,
		Body:         "",
	}

	actualDoc := newDoc("000000000020", 450)

	if !reflect.DeepEqual(expectedDoc, actualDoc) {
		t.Errorf("expected: %v, got: %v", expectedDoc, actualDoc)
	}
}
