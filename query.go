package qb

// Field lookup by document ID (key)
func q1(keySpace int64) *QueryPayload {
	_, key := existingKey(prefix, keySpace)

	return &QueryPayload{
		QueryType:  "Q1",
		Projection: []string{"address"},
		Selection: []Filter{
			{"_id", key, false},
		},
	}
}

// Unique lookup by document field
func q2(keySpace int64) *QueryPayload {
	i, key := existingKey(prefix, keySpace)
	alphabet := newAlphabet(i, key)

	return &QueryPayload{
		QueryType:  "Q2",
		Projection: []string{"address"},
		Selection: []Filter{
			{"email", newEmail(alphabet), false},
		},
	}
}

// Range search by document field
func q3(keySpace int64) *QueryPayload {
	i, _ := existingKey(prefix, keySpace)

	return &QueryPayload{
		QueryType:  "Q3",
		Projection: []string{"address"},
		Selection: []Filter{
			{"localgroup", newGroup(i, 100), false},
		},
	}
}

// Composite search using equality predicate and text search
func q4(keySpace int64) *QueryPayload {
	i, _ := existingKey(prefix, keySpace)

	return &QueryPayload{
		QueryType:  "Q4",
		Projection: []string{"firstname", "lastname"},
		Selection: []Filter{
			{"address.zip", newZipCode(i), false},
			{"address.street", "z" + newGroup(i, 1000*(1+i%3)), true},
		},
	}
}

// Document lookup by document ID (key)
func q5(keySpace int64) *QueryPayload {
	_, key := existingKey(prefix, keySpace)

	return &QueryPayload{
		QueryType:  "Q5",
		Projection: []string{},
		Selection: []Filter{
			{"_id", key, false},
		},
	}
}
