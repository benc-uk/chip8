# Common variables
VERSION := 0.0.1
BUILD_INFO := Manual build 

SRC_DIR := cmd

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go

.PHONY: help run lint lint-fix
.DEFAULT_GOAL := help

help: ## This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## Lint & format, will not fix but sets exit code on error
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SRC_DIR); $(GOLINT_PATH) run --modules-download-mode=mod *.go --out-format=github-actions

lint-fix: ## Lint & format, will try to fix errors and modify code
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SRC_DIR); golangci-lint run --modules-download-mode=mod *.go --fix

build: ## Run a local build without a container
	@mkdir -p bin
	go mod tidy
	go build -o bin/chip8 ./$(SRC_DIR)/...

run: ## Run application, used for local development
	air -c .air.toml