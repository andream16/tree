PROJECT?=advisor

default: build

# Build
build: 
	go build -o bin/$(PROJECT) -v

# Test
test: go test ./...
test-verbose: go test -v ./...
test-cover: go cover ./...
