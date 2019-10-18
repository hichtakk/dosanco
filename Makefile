NAME := dosanco
REVISION := $(shell git rev-parse --shrt HEAD)

RELEASE_DIR=build
GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(THIS_GOOS)
GOARCH=$(THIS_GOARCH)

.DEFAULT_GOAL := help
BUILD_TARGETS= \
	build-linux-amd64 \
	build-linux-arm \
	build-darwin-amd64

.PHONY: help clean build all fmt lint $(BUILD_TARGETS)

fmt: ## format
	go fmt

lint: ## Examine source code and lint
	go vet ./...
	golint -set_exit_status ./...

all: $(BUILD_TARGETS)

build: $(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH) $(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH) ## build dosanco and dosanco-apiserver

#linux: ## build AMD64 linux binary
#	GOOS=linux GOARCH=amd64 go build -o build/dosanco_linux_amd64 cli/main.go
#	#see https://qiita.com/keijidosha/items/5f4a68a3341a44a25ab9
#	#GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o build/linux/amd64/dosanco-apiserver main.go
#	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -v -o build/dosanco-apiserver_linux_amd64 main.go

build-linux-amd64: ## build AMD64 linux binary
	#@$(MAKE) build GOOS=linux GOARCH=amd64 LDFLAGS="-ldflags '-linkmode external -extldflags \"-static\"'"
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64: ## build ARM64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=amd64

$(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH): ## Build dosanco command-line client
	@echo "==> Build dosanco for ${GOOS}-${GOARCH}"
	@GO111MODULE=on go build -a -v -o $(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH) cli/main.go

$(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH): ## Build dosanco api server
	@echo '==> Build dosanco-apiserver for ${GOOS}-${GOARCH}'
	@GO111MODULE=on CGO_ENABLED=1 go build -a -v $(LDFLAGS) -o $(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH) main.go

update-package: ## Update dependency packages
	@go mod tidy
	@go get -u

clean: ## Clean up built files
	@echo '==> Remove built files ./build/...'
	@rm -rf build/*

help: ## Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort