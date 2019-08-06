all: build run

build:
	go build -o bin/jump main.go

run:
	./bin/jump

test:
	go test -cover -v

.PHONY: build run test
