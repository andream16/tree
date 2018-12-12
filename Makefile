PROJECT?=adviser
BINPATH?= bin/$(PROJECT)

default: 
	build
		./$(BINPATH)

build: 
	go build -o $(BINPATH) -v

# Test
test: go test ./...
test-verbose: go test -v ./...
test-cover: go test -cover ./...
