SHELL := /usr/bin/env bash

BINARY ?= depctl
BUILD_DIR ?= dist

.PHONY: help fmt fmt-check vet test race build smoke verify clean

help:
	@printf '%s\n' \
		'Targets:' \
		'  fmt        Format Go code' \
		'  fmt-check  Fail if Go code is not formatted' \
		'  vet        Run go vet' \
		'  test       Run unit and integration tests' \
		'  race       Run tests with the race detector' \
		'  build      Build the CLI into dist/' \
		'  smoke      Run CLI smoke tests against fixtures' \
		'  verify     Run the full local quality gate' \
		'  clean      Remove build output'

fmt:
	gofmt -w $$(find . -path './.git' -prune -o -name '*.go' -print)

fmt-check:
	@test -z "$$(gofmt -l $$(find . -path './.git' -prune -o -name '*.go' -print))"

vet:
	go vet ./...

test:
	go test ./...

race:
	go test -race ./...

build:
	mkdir -p "$(BUILD_DIR)"
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$(BUILD_DIR)/$(BINARY)" .

smoke:
	bash scripts/smoke.sh

verify: fmt-check vet test race build smoke

clean:
	rm -rf "$(BUILD_DIR)"
