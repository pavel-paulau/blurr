package qb

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewKey(t *testing.T) {
	prefix := "prefix"
	i := int64(123456789)
	actualKey := newKey(prefix, i)
	expectedKey := "prefix-3626229738a3b9"

	if expectedKey != actualKey {
		t.Errorf("expected: %v, got: %v", expectedKey, actualKey)
	}
}

func TestExistingKey(t *testing.T) {
	prefix := "prefix"
	i, actualKey := existingKey(prefix, 1000000)
	expectedKey := newKey(prefix, i)

	if expectedKey != actualKey {
		t.Errorf("expected: %v, got: %v", expectedKey, actualKey)
	}
}

func TestDoc(t *testing.T) {
	expectedDoc := Doc{
		FirstName: "ckyK2nI3",
		LastName:  "v1JDXLP0",
		Email:     "xqo6snEY@Q7tw3G.com",
		Address: address{
			City:      "P5dV9ajx",
			County:    "OYFCTKA",
			Country:   "H7Wdlg8Xr",
			FullState: "Nebraska",
			State:     "NE",
			Street:    "1789 bc614e 1e240 Place",
			Zip:       86789,
		},
		Category:    4,
		Balance:     163.06,
		DateOfBirth: "2007-02-16T00:00:00Z",
		Avatar:      "https://www.gravatar.com/avatar/dV9ajxOYFCTKAH7Wdlg8Xrf9eCmzynpv",
		Company:     "f9eCmzynpv inc.",
		Age:         23,
		LocalGroup:  "12d687",
	}
	actualDoc := newDoc(123456789, "prefix-3626229738a3b9", 0)

	if !reflect.DeepEqual(expectedDoc, actualDoc) {
		t.Errorf("expected: %+v, got: %+v", expectedDoc, actualDoc)
	}
}

func TestDocSize(t *testing.T) {
	for _, size := range []int64{512, 768, 1024, 2048} {
		doc := newDoc(123456789, "prefix-003626229738a3b9", size)

		b, err := json.Marshal(&doc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(b) != int(size) {
			t.Errorf("expected %v, got %v", size, len(b))
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TestUniqueness(t *testing.T) {
	values := []string{}
	for i := 0; i < 1e4; i++ {
		key := newKey("prefix", int64(i))
		s := newString(int64(i), key, 64)
		if stringInSlice(s, values) {
			t.Fatalf("Duplicate %d - %s", i, s)
		}
		values = append(values, s)
	}
}

func BenchmarkNewString(b *testing.B) {
	i := int64(123456789)
	k := "prefix-003626229738a3b9"
	for n := 0; n < b.N; n++ {
		newString(i, k, 64)
	}
}

func BenchmarkNewAlphabet(b *testing.B) {
	k := "prefix-003626229738a3b9"
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newAlphabet(i, k)
	}
}

func BenchmarkNewName(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newFirstName(a)
	}
}

func BenchmarkNewEmail(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newEmail(a)
	}
}

func BenchmarkNewCategory(b *testing.B) {
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newCategory(i)
	}
}

func BenchmarkNewDOB(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newDateOfBirth(a)
	}
}

func BenchmarkNewCompany(b *testing.B) {
	i := int64(123456789)
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newCompany(a, i)
	}
}

func BenchmarkNewAge(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newAge(a)
	}
}

func BenchmarkNewState(b *testing.B) {
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newState(i)
	}
}

func BenchmarkNewStreet(b *testing.B) {
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newStreet(i)
	}
}

func BenchmarkNewAvatar(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newAvatar(a)
	}
}

func BenchmarkNewBalance(b *testing.B) {
	a := newAlphabet(123456789, "prefix-003626229738a3b9")
	for n := 0; n < b.N; n++ {
		newBalance(a)
	}
}

func BenchmarkNewGroup(b *testing.B) {
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newGroup(i, 10)
	}
}

func BenchmarkNewKey(b *testing.B) {
	prefix := "prefix"
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		newKey(prefix, i)
	}
}

func BenchmarkExistingKey(b *testing.B) {
	prefix := "prefix"
	i := int64(123456789)
	for n := 0; n < b.N; n++ {
		existingKey(prefix, i)
	}
}

func BenchmarkBaseDoc(b *testing.B) {
	i := int64(123456789)
	var size int64
	k := "prefix-003626229738a3b9"
	for n := 0; n < b.N; n++ {
		newDoc(i, k, size)
	}
}

func BenchmarkNewDoc512(b *testing.B) {
	i := int64(123456789)
	size := int64(512)
	k := "prefix-003626229738a3b9"
	for n := 0; n < b.N; n++ {
		newDoc(i, k, size)
	}
}

func BenchmarkNewDoc1024(b *testing.B) {
	i := int64(123456789)
	size := int64(1024)
	k := "prefix-003626229738a3b9"
	for n := 0; n < b.N; n++ {
		newDoc(i, k, size)
	}
}

func BenchmarkNewDoc2048(b *testing.B) {
	i := int64(123456789)
	size := int64(2048)
	k := "prefix-003626229738a3b9"
	for n := 0; n < b.N; n++ {
		newDoc(i, k, size)
	}
}
