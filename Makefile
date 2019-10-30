NAME := dosanco

RELEASE_DIR=build
GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(THIS_GOOS)
GOARCH=$(THIS_GOARCH)
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const version' main.go)))
REVISION=$(shell git rev-parse HEAD)

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

build-linux-amd64: ## build AMD64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64: ## build ARM64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=amd64

$(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH): ## Build dosanco command-line client
	@echo "==> Build dosanco for ${GOOS}-${GOARCH}"
	echo ${REVISION}
	@GO111MODULE=on go build -ldflags "-X github.com/hichikaw/dosanco/cmd.revision=${REVISION}" -a -v -o $(RELEASE_DIR)/dosanco_$(GOOS)_$(GOARCH) cli/main.go

$(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH): ## Build dosanco api server
	@echo '==> Build dosanco-apiserver for ${GOOS}-${GOARCH}'
	@GO111MODULE=on CGO_ENABLED=1 go build -ldflags "-X main.revision=${REVISION}"  -a -v $(LDFLAGS) -o $(RELEASE_DIR)/dosanco-apiserver_$(GOOS)_$(GOARCH) main.go

docker-amd64: build-linux-amd64 ## build docker image for AMD64 architecture
	@docker build -t docker.pkg.github.com/hichikaw/dosanco:${VERSION} .

update-package: ## Update dependency packages
	@go mod tidy
	@go get -u

clean: ## Clean up built files
	@echo '==> Remove built files ./build/...'
	@rm -rf build/*

help: ## Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort