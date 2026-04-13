BINARY_NAME := qmai
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

PUBLIC_REPO ?= /Users/feixiaoliang/Documents/www/qmai-cli-public
RELEASE_VERSION ?=
PUBLIC_VERSION ?=
PUBLIC_TAG ?=
SOURCE_REF ?= HEAD
# 若设置（如 0.1.1），publish-open-source 成功后会顺带执行 build-release-archives
ARCHIVE_VERSION ?=
RELEASE_ARGS := $(if $(strip $(RELEASE_VERSION)),--release-version $(strip $(RELEASE_VERSION)),)
ARCHIVE_ARGS := $(if $(strip $(ARCHIVE_VERSION)),--build-archives $(strip $(ARCHIVE_VERSION)),)

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
	./scripts/publish-open-source.sh --public-repo "$(PUBLIC_REPO)" --source-ref "$(SOURCE_REF)" $(RELEASE_ARGS) --version "$(PUBLIC_VERSION)" $(ARCHIVE_ARGS) --dry-run

publish-open-source:
	./scripts/publish-open-source.sh --public-repo "$(PUBLIC_REPO)" --source-ref "$(SOURCE_REF)" $(RELEASE_ARGS) --version "$(PUBLIC_VERSION)" --tag "$(PUBLIC_TAG)" $(ARCHIVE_ARGS)

init-open-source-repo:
	./scripts/init-open-source-repo.sh --public-repo "$(PUBLIC_REPO)"

all: fmt vet test build
