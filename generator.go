package qb

type kvPayload struct {
	key   string
	value *doc
}

const prefix = "user-profile"

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func generatePayload(workerID, numDocs, docSize int64) chan kvPayload {
	payload := make(chan kvPayload, min(1e3, numDocs))

	go func() {
		defer close(payload)

		for i := int64(0); i < numDocs; i++ {
			j := workerID*numDocs + i
			key := newKey(prefix, j)
			doc := newDoc(j, key, docSize)
			payload <- kvPayload{key, &doc}
		}
	}()

	return payload
}
