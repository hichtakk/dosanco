.PHONY: fmt lint build clean

fmt:
	go fmt

lint:
	golint
	golint handler
	golint config
	golint db
	golint model
	golint cmd

build: build/dosanco-apiserver build/dosanco

build/dosanco-apiserver:
	go build -o build/dosanco-apiserver main.go

build/dosanco:
	go build -o build/dosanco cli/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/dosanco cli/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o build/linux/amd64/dosanco-apiserver main.go

clean:
	rm -rf build/*