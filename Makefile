# Binary dependencies
golangci-lint := go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest
goreleaser := go run github.com/goreleaser/goreleaser/v2@latest

test:
	go test -race -v ./...

build: download-deps tidy-deps compile test lint

fmt:
	gofmt -l -s -w .

download-deps:
	go mod download

tidy-deps:
	go mod tidy

update-deps:
	go get -u -t ./...
	go mod tidy

compile:
	go build -v ./...

lint:
	$(golangci-lint) run
	$(goreleaser) check

release:
	$(goreleaser) release --clean