.DEFAULT_GOAL := run

fmt:
	go fmt .
	go fmt ./sqlite
.PHONY:fmt

lint: fmt
	golint .
	golint ./sqlite
.PHONY:lint

vet: lint
	go vet .
	go vet ./sqlite
.PHONY:vet

run: vet
	go run .
.PHONY:run

build: vet
	go build .
.PHONY:build
