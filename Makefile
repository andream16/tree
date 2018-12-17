default: build

workdir:
	mkdir -p bin

build: bin/tree

bin/tree:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tree .

test: test-all

test-all:
	go test -v -cover -race ./...