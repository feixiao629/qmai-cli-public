BINARY_NAME := qmai
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

PUBLIC_REPO ?=
PUBLIC_VERSION ?=
PUBLIC_TAG ?=
SOURCE_REF ?= HEAD

.PHONY: build test lint clean install fmt vet publish-open-source-dry-run publish-open-source init-open-source-repo

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

publish-open-source-dry-run:
	./scripts/publish-open-source.sh --public-repo "$(PUBLIC_REPO)" --source-ref "$(SOURCE_REF)" --version "$(PUBLIC_VERSION)" --dry-run

publish-open-source:
	./scripts/publish-open-source.sh --public-repo "$(PUBLIC_REPO)" --source-ref "$(SOURCE_REF)" --version "$(PUBLIC_VERSION)" --tag "$(PUBLIC_TAG)"

init-open-source-repo:
	./scripts/init-open-source-repo.sh --public-repo "$(PUBLIC_REPO)"

all: fmt vet test build
