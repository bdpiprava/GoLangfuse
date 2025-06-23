.PHONY: tests

# To run unit tests
tests:
	@echo "Running tests"
	go test -count=1 -race ./...

build:
	@echo "Building the application"
	go build -o golangfuse ./main