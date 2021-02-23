NAME := dosanco
.DEFAULT_GOAL := help

RELEASE_DIR := build
BUILD_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-amd64

GOVERSION = $(shell go version)
THIS_GOOS = $(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH = $(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS = $(THIS_GOOS)
GOARCH = $(THIS_GOARCH)
VERSION = $(patsubst "%",%,$(lastword $(shell grep 'const version' main.go)))
REVISION = $(shell git rev-parse HEAD)

.PHONY: help clean update-package docker-linux-amd64 build all fmt lint $(BUILD_TARGETS)

fmt: ## format
	go fmt

lint: ## Examine source code and lint
	go vet ./...
	golint -set_exit_status ./...

all: $(BUILD_TARGETS) ## build for all platform

build: $(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH) $(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH) ## build dosanco and dosanco-apiserver

build-linux-amd64: ## build AMD64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64: ## build ARM64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=amd64

$(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH): ## Build dosanco command-line client
	@printf "\e[32m"
	@echo "==> Build dosanco for ${GOOS}-${GOARCH}"
	@printf "\e[90m"
	@GO111MODULE=on go build -ldflags "-X github.com/hichtakk/dosanco/cmd.revision=${REVISION}" -a -v -o $(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH) cli/dosanco/main.go
	@printf "\e[m"

$(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH): ## Build dosanco api server
	@printf "\e[32m"
	@echo '==> Build dosanco-apiserver for ${GOOS}-${GOARCH}'
	@printf "\e[90m"
	@GO111MODULE=on CGO_ENABLED=1 go build -ldflags "-X main.revision=${REVISION}"  -a -v $(LDFLAGS) -o $(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH) main.go route.go
	@printf "\e[m"

release-github: ## tag and release to github
	@ghr ${VERSION} build

docker-linux-amd64: build-linux-amd64 ## build docker image for AMD64 architecture
	@docker build -t hichtakk/dosanco:${VERSION} .

update-package: ## Update dependency packages
	@go mod tidy
	@go get -u

clean: ## Clean up built files
	@printf "\e[32m"
	@echo '==> Remove built files ./build/...'
	@printf "\e[90m"
	@ls -1 ./build
	@rm -rf build/*
	@printf "\e[m"

help: ## Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort

echo:
	@#@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST)
	@grep -E '^.+:.* +##' $(MAKEFILE_LIST) | sed -e 's/\(##.*\)$$/"\x1b[31m" \1 "\x1b[0m"/g'
	echo $(BUILD_TARGETS)

