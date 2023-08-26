.DEFAULT_GOAL := run

fmt:
	go fmt ./cmd/http
	go fmt ./sqlite
.PHONY:fmt

lint: fmt
	golint ./cmd/http
	golint ./sqlite
.PHONY:lint

vet: lint
	go vet ./cmd/http
	go vet ./sqlite
.PHONY:vet

run: vet
	go run ./cmd/http
.PHONY:run

build: vet
	go build ./cmd/http
.PHONY:build
