BINARY_NAME := qmai
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build test lint clean install fmt vet

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

install:
	go install $(LDFLAGS) .

test:
	go test ./... -v

test-short:
	go test ./... -short

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -rf bin/

all: fmt vet test build
