.PHONY: build test lint fmt ci

# Build (no-op for library)
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Check formatting
fmt:
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)

# Run linter
lint:
	golangci-lint run

# Run all CI checks
ci: fmt lint test
