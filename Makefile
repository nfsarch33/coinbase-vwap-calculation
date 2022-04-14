CompiledFileName=build/vwap-cli # Name of compiled binary
TestInclusion=$(shell go list ./... | grep -Ewv "test|cmd|dummy")

# Default environment variables
env=local
scenario=all

install:
	@go mod tidy

clean:
	@rm -f $(CompiledFileName)

build:
	@go build -o $(CompiledFileName) cmd/vwap/main.go

run:
	@go run cmd/vwap/main.go

test:
	@go test -v -race -timeout 10000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.raw.out $(TestInclusion)

lint_ci:
	@golangci-lint run ./...

linter_install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
