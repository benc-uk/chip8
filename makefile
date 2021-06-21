# Common variables
VERSION := 1.0.0
BUILD_INFO := Manual build 

SRC_DIR := cmd
GO_PKG := github.com/benc-uk/chip8

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go
GOOS ?= linux

.PHONY: help run lint lint-fix test run
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 🔍 Lint & format, will not fix but sets exit code on error
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run --modules-download-mode=mod --out-format=github-actions $(SRC_DIR)/...

lint-fix: ## 📝 Lint & format, will try to fix errors and modify code
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	golangci-lint run --modules-download-mode=mod --fix $(SRC_DIR)/...

build: ## 🔨 Run a local build without a container
	@mkdir -p bin
	go mod tidy
	GOOS=linux go build -o bin/chip8 ./$(SRC_DIR)/chip8/...
	GOOS=windows go build -o bin/chip8.exe ./$(SRC_DIR)/chip8/...
	GOOS=js GOARCH=wasm go build -o web/chip8.wasm ./$(SRC_DIR)/chip8wasm/...

run: ## 🏃‍ Run application, used for local development
	DISPLAY=192.168.0.34:0 air -c .air.toml

test: ## 🤡 Run those sweet unit tests to give the illusion of testing
	go test -v -count 1 $(GO_PKG)/pkg/chip8/... 