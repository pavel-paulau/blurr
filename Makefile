build:
	go build -v ./cmd/cbload
	go build -v ./cmd/cbrun
	go build -v ./cmd/elasticload
	go build -v ./cmd/elasticrun
	go build -v ./cmd/ftsrun
	go build -v ./cmd/mongoload
	go build -v ./cmd/mongorun

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test -v -cover -race -coverprofile=coverage.out

bench:
	go test -v -run=AAA -test.benchmem -bench=.

clean:
	rm -f coverage.out cbload cbrun elasticload elasticrun ftsrun mongoload mongorun
