build:
	go build -v -o cbl   ./cbload
	go build -v -o cbr   ./cbrun
	go build -v -o esl   ./elasticload
	go build -v -o esr   ./elasticrun
	go build -v -o ftsr  ./ftsrun
	go build -v -o mgl   ./mongoload
	go build -v -o mgr   ./mongorun

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test -v -cover -race -coverprofile=coverage.out

bench:
	go test -v -run=AAA -test.benchmem -bench=.

clean:
	rm -f coverage.out cbl cbr esl esr ftsr mgl mgr
