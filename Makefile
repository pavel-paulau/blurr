build:
	go build -v

fmt:
	gofmt -w -s *.go

test:
	go test -v -cover -race .

bench:
	go test -v -run=AAA -test.benchmem -bench=.
