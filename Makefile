.DEFAULT_GOAL := run

http_path := ./cmd/http
sql_pkg := ./dblog
bins := ./http

fmt:
	find . -type f -name '*.go' | xargs -I{} go fmt {}
.PHONY:fmt

lint: fmt
	golangci-lint run || true
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
