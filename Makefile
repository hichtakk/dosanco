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

clean:
	rm build/*