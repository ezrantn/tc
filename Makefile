BINARY_NAME=trc
PKG=.
OUTPUT_DIR=bin
EXECUTABLE_FILE=cmd/main.go

GO=go
GOTEST=$(GO) test -v
GOBUILD=$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(EXECUTABLE_FILE)
GOLINT=golangci-lint run
GOTIDY=$(GO) mod tidy
GOCOVOPENHTML=$(GO) tool cover -html=coverage.out
GOCOV=$(GO) test . -coverprofile=coverage.out
GOBENCH=$(GO) test -bench=. -benchtime=10s ./benches

default: build

.PHONY: build
build:
	@echo "Building the application..."
	$(GOBUILD)

.PHONY: run
run: build
	@echo "Running the application..."
	$(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) $(PKG)

.PHONY: lint
lint:
	@echo "Running linter..."
	$(GOLINT)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR) coverage.out 

.PHONY: tidy
tidy:
	@echo "Tidying up Go modules..."
	$(GOTIDY)

.PHONY: open-html
open-html:
	@echo "Open Code Coverage in HTML"
	$(GOCOVOPENHTML)

.PHONY: cov
cov: coverage open-html
	@echo "Test coverage generated and opened."

.PHONY: coverage
coverage:
	@echo "Generating coverage report..."
	$(GOCOV)

.PHONY: bench
bench: 
	@echo "Benchmarking..."
	$(GOBENCH)

.PHONY: all
all: tidy lint test build 
