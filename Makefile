.PHONY: tests

# Packages
PKG := $(shell go list ./... | grep -v 'github.com/bdpiprava/GoLangfuse/test-integration')

test-deps:
	@go install gotest.tools/gotestsum@v1.12.1

# To run unit tests
tests:
	@echo "Running all tests"
	@make test-unit
	@make test-integration

test-unit: test-deps
	@echo "Running unit tests"
	@gotestsum --format=testname -- -v -coverprofile=coverage_unit.txt -race -cover -covermode=atomic $(PKG) -short

test-integration: test-deps
	@echo "Running integration tests"
	@gotestsum --format testname -- -v -coverprofile=coverage_integration.txt -race -cover -covermode=atomic -coverpkg=./... -run "^*.IntegrationTestSuite$$" ./...

build:
	@echo "Building the application"
	go build -o golangfuse ./main