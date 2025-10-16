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
	@echo "✓ Upgrading all compile and test dependencies..."
	@go get -u -t ./...
	@echo "✓ Upgrading all dev dependencies..."
	@go get tool
	@echo "✓ Tidying up the go.mod file..."
	@go mod tidy

compile:
	go build -v ./...

lint:
	go tool golangci-lint run
	go tool goreleaser check

release:
	go tool goreleaser release --clean