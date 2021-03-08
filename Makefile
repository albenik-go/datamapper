SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: help init update protoc generate lint test test-build

help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Prepare project for development
	git config core.hooksPath .githooks
	go mod tidy
	go generate ./...

update: ## Update go.mod dependencies
	go get -u ./...
	go mod tidy

generate: protoc ## Run code generation
	go generate ./...

lint: ## Run golangci-lint
	go mod tidy
	golangci-lint run

test: ## Run only unit tests
	go test ./...
