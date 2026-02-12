# List of effective go files
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" | egrep -v "^\./\.go" | grep -v _test.go)

# List of packages except testsutils
PACKAGES ?= $(shell go list ./... | grep -v "mock" )

# Build folder
BUILD_FOLDER = build

# Test coverage variables
COVERAGE_BUILD_FOLDER = $(BUILD_FOLDER)/coverage

UNIT_COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/ut_cov.out
UNIT_COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/ut_index.html

INTEGRATION_COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/it_cov.out
INTEGRATION_COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/it_index.html

# Test lint variables
GOLANGCI_VERSION = v2.8.0

MOCKGEN = go run github.com/golang/mock/mockgen@v1.6.0 \

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build/coverage:
	@mkdir -p build/coverage

unit-test: build/coverage
	@go test -covermode=count -coverprofile $(UNIT_COVERAGE_OUT) $(PACKAGES)

unit-test-cov: unit-test
	@go tool cover -html=$(UNIT_COVERAGE_OUT) -o $(UNIT_COVERAGE_HTML)

fix-lint: lint-fix ## Run linter to fix issues
	

# @misspell -error $(GOFILES)
test-lint: lint ## Check linting

lint:
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION) run

lint-fix:
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION) run --fix --issues-exit-code=0

.PHONY: lint lint-fix

integration-test: build/coverage
	@go test -covermode=count -coverprofile $(INTEGRATION_COVERAGE_OUT) -v --tags integration ${PACKAGES}

integration-test-cov: integration-test
	@go tool cover -html=$(INTEGRATION_COVERAGE_OUT) -o $(INTEGRATION_COVERAGE_HTML)

mocks:
	${MOCKGEN} -source ethereum/execution/client/client.go -destination ethereum/execution/client/mock/client.go -package mock Client
	${MOCKGEN} -source ethereum/consensus/client/client.go -destination ethereum/consensus/client/mock/client.go -package mock Client
	${MOCKGEN} -source net/jsonrpc/client.go -destination net/jsonrpc/testutils/client.go -package testutils Client
