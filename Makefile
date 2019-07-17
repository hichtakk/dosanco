fmt:
	go fmt

lint:
	golint
	golint handler
	golint config
	golint db
	golint model
	golint cmd

dctl:
	go build -o dctl cli/dctl.go