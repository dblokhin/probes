GO ?= go
GOTEST = $(GO) test
GOLANGCI_LINT_VERSION := v1.60.2
BIN_DIR := $(shell go env GOPATH)/bin

test:
	$(GOTEST) -cover -count=1 ./...

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...
	$(BIN_DIR)/golangci-lint run ./...

deps:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) $(GOLANGCI_LINT_VERSION)
	@go install github.com/mattn/goveralls@latest

cover:
	goveralls

.PHONY: test fmt lint deps cover
