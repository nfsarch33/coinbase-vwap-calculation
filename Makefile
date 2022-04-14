CompiledFileName=build/vwap-cli # Name of compiled binary
TestInclusion=$(shell go list ./... | grep -Ewv "test|cmd|dummy")

# Default environment variables
env=local
scenario=all

install:
	@go mod tidy
	@go install -v ./...

build:
	@go build -o $(CompiledFileName) cmd/vwap/main.go

clean:
	@rm -rf $(CompiledFileName)

start:
	@./$(CompiledFileName) -verbose -wsurl "wss://ws-feed.exchange.coinbase.com" -window-size 200 -pairs "BTC-USD,ETH-USD,ETH-BTC"

run:
	@go run cmd/vwap/main.go

test:
	@go test -v -race -tags=${scenario} -timeout 10000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.raw.out $(TestInclusion)

test_integration:
	@go test -v -race -tags=integration -timeout 10000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.raw.out $(TestInclusion)

lint:
	@golangci-lint run ./...

linter_install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
