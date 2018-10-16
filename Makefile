all: build run

build:
	go build -o bin/jump ./cmd

run:
	./bin/jump

test:
	go test -cover -v

.PHONY: build run test
