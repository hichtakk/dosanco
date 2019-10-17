.PHONY: fmt lint build clean

NAME := dosanco
#VERSION := $(gobump show -r)
REVISION := $(shell git rev-parse --shrt HEAD)
LDFLAGS := "-linkmode external -extldflags "-static""

.DEFAULT_GOAL := help

fmt: ## format
	go fmt

lint: ## Examine source code and lint
	go vet ./...
	golint -set_exit_status ./...

build: linux macos ## build 

build/dosanco-apiserver: ## Build dosanco api server
	go build -o build/dosanco-apiserver main.go

build/dosanco: ## Build dosanco command-line client
	go build -o build/dosanco cli/main.go

linux: ## build AMD64 linux binary
	GOOS=linux GOARCH=amd64 go build -o build/dosanco_linux_amd64 cli/main.go
	#see https://qiita.com/keijidosha/items/5f4a68a3341a44a25ab9
	#GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o build/linux/amd64/dosanco-apiserver main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -v -o build/dosanco-apiserver_linux_amd64 main.go

macos: ## build AMD64 darwin binary
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -a -v -o build/dosanco-apiserver/dosanco-apiserver_darwin_amd64 main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -a -v -o build/dosanco/dosanco_darwin_amd64 cli/main.go

update-package: ## Update dependency packages
	@go mod tidy
	@go get -u

clean: ## Clean up built files
	@echo '==> Remove built files ./build/...'
	@echo ''
	rm -rf build/*

help: ## Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort