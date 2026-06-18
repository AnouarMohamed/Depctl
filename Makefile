SHELL := /usr/bin/env bash

BINARY ?= depctl
BUILD_DIR ?= dist

.PHONY: help fmt fmt-check vet test race build smoke verify release clean

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
		'  release    Create a guarded release tag, VERSION=vX.Y.Z' \
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
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$(BUILD_DIR)/$(BINARY)" ./cmd/depctl

smoke:
	bash scripts/smoke.sh

verify: fmt-check vet test race build smoke

release:
	@test -n "$(VERSION)" || (echo 'VERSION is required, for example: make release VERSION=v1.1.0' >&2; exit 2)
	bash scripts/release.sh "$(VERSION)"

clean:
	rm -rf "$(BUILD_DIR)"
