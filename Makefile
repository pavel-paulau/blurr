build:
	go build -v

fmt:
	gofmt -w -s *.go

test:
	go test -v -cover -race .

clean:
	rm -fr nb
