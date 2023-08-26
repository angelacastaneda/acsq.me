.DEFAULT_GOAL := run

fmt:
	go fmt .
.PHONY:fmt

lint: fmt
	golint .
.PHONY:lint

vet: lint
	go vet .
.PHONY:vet

run: vet
	go run .
.PHONY:run

build: vet
	go build .
.PHONY:build
