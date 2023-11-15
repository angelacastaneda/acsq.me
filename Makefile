.DEFAULT_GOAL := run

http_path := ./cmd/http
sql_pkg := ./dblog
bins := ./http

fmt:
	go fmt $(http_path)
	go fmt $(sql_pkg)
.PHONY:fmt

lint: fmt
	golint $(http_path)
	golint $(sql_pkg)
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
