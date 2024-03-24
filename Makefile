.DEFAULT_GOAL := run

http_path := ./cmd/http
sql_pkg := ./dblog
bins := ./http

fmt:
	go fmt $(http_path)
	go fmt $(sql_pkg)
.PHONY:fmt

lint: fmt
	golangci-lint run $(http_path) || true
	golangci-lint run $(sql_pkg) || true
.PHONY:lint

vet: lint
	go vet $(http_path)
	go vet $(sql_pkg)
.PHONY:vet

run: vet
	go run $(http_path)
.PHONY:run

build:
	go build $(http_path)
.PHONY:build

clean:
	rm -fv $(bins)
.PHONY:clean
