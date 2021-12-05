.DEFAULT_GOAL := build

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

MAIN_APP="bin/generics"

.PHONY: fmt vet build

all: test vet fmt build

list:
	@cat $(MAKEFILE_LIST) | awk -F: '/^([a-zA-Z_-]+):/ {print $$1}'

reset:
	go clean -testcache ./...

clean:
	rm -f $(MAIN_APP)

test:
	go test ./...

retest: reset test

vet:
	go vet ./...

shadow:
	shadow ./...

fmt:
	go fmt ./...

lint:
	golint ./...

build:
	go clean
	go build -o $(MAIN_APP) main.go

run:
	go run main.go

