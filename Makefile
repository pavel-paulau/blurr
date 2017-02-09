build:
	go build -v -o cbl ./cbload
	go build -v -o mgl ./mongoload

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test -v -cover -race .

bench:
	go test -v -run=AAA -test.benchmem -bench=.
