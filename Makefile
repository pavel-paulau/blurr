build:
	go build -v -o cb_load ./cbload

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test -v -cover -race .

bench:
	go test -v -run=AAA -test.benchmem -bench=.
