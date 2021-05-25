SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: help
help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Prepare project for development
	git config core.hooksPath .githooks
	go mod tidy
	go generate ./...

.PHONY: update
update: ## Update go.mod dependencies
	go get -u ./...
	go mod tidy
	go generate ./...

.PHONY: generate
generate: ## Run code generation
	go generate ./...

.PHONY: lint
lint: ## Run golangci-lint
	go mod tidy
	golangci-lint run

.PHONY: test
test: ## Run only unit tests
	go test ./...
