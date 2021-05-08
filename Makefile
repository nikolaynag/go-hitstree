.PHONY: all gofmt test clean run-example

all: test

gofmt:
	go fmt ./...

test:
	go test -cover -v ./...

clean:
	rm -rf build/

run-example:
	go build -o build/bin/example example/example.go
	./build/bin/example
