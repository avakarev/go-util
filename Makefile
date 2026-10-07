lint-install:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

lint:
	@echo ">> Running golangci-lint..."
	@golangci-lint version
	@golangci-lint run ./...

vet:
	@echo ">> Vetting..."
	@go vet ./...

test:
	@echo ">> Running tests..."
	@go test -v -race ./...
.PHONY: test

ci: lint vet test
