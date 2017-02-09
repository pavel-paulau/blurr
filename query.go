package qb

func q2(keySpace int64) queryPayload {
	i, key := existingKey(prefix, keySpace)
	alphabet := newAlphabet(i, key)
	arg := newEmail(alphabet)
	return queryPayload{"email", arg}
}

func q3(keySpace int64) queryPayload {
	i, _ := existingKey(prefix, keySpace)
	arg := newGroup(i, 100)
	return queryPayload{"localgroup", arg}
}
