.PHONY: fmt lint build clean

NAME := dosanco
#VERSION := $(gobump show -r)
REVISION := $(shell git rev-parse --shrt HEAD)
LDFLAGS := "-linkmode external -extldflags "-static""


fmt:
	go fmt

lint:
	go vet ./...
	golint -set_exit_status ./...

build: build/dosanco-apiserver build/dosanco

build/dosanco-apiserver:
	go build -o build/dosanco-apiserver main.go

build/dosanco:
	go build -o build/dosanco cli/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/dosanco cli/main.go
	#see https://qiita.com/keijidosha/items/5f4a68a3341a44a25ab9
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o build/linux/amd64/dosanco-apiserver main.go

update-package:
	@go mod tidy
	@go get -u

clean:
	rm -rf build/*