# Define some VCS context
PARENT_BRANCH ?= master

# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_DIR_NAME ?= _report
GO_TEST_COVERAGE_FILE_NAME ?= ${GO_TEST_COVERAGE_DIR_NAME}/coverage.out
GO_TEST_COVERAGE_HTML_FILE_NAME ?= ${GO_TEST_COVERAGE_DIR_NAME}/coverage.html

# Set flags for `gofmt`
GOFMT_FLAGS ?= -s

# Set a default `min_confidence` value for `golint`
GOLINT_MIN_CONFIDENCE ?= 0.1


all: install-deps build install

clean:
	go clean -i -x ./...

build:
	go build -v ./...

install:
	go install ./...

install-deps:
	go get -d ./...

install-deps-dev: install-deps
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

update-deps:
	go get -d -u ./...

update-deps-dev: update-deps
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

test:
	go test -v ./...

test-with-coverage:
	go test -cover ./...

test-with-coverage-json-output:
	go test -cover -json ./...

test-with-coverage-profile ${GO_TEST_COVERAGE_FILE_NAME}:
	@mkdir -p "$$(dirname "${GO_TEST_COVERAGE_FILE_NAME}")"
	go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile ${GO_TEST_COVERAGE_FILE_NAME} ./...

generate-test-coverage-html ${GO_TEST_COVERAGE_HTML_FILE_NAME}:
	@mkdir -p "$$(dirname "${GO_TEST_COVERAGE_HTML_FILE_NAME}")"
	go tool cover -html=${GO_TEST_COVERAGE_FILE_NAME} -o ${GO_TEST_COVERAGE_HTML_FILE_NAME}

test-with-coverage-profile-html: test-with-coverage-profile generate-test-coverage-html
${GO_TEST_COVERAGE_DIR_NAME}: ${GO_TEST_COVERAGE_FILE_NAME} ${GO_TEST_COVERAGE_HTML_FILE_NAME}

format-lint:
	errors=$$(gofmt -l ${GOFMT_FLAGS} .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

import-lint:
	errors=$$(goimports -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

style-lint:
	errors=$$(golint -min_confidence=${GOLINT_MIN_CONFIDENCE} ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

lint: install-deps-dev format-lint import-lint style-lint

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	goimports -w .

vet:
	go vet ./...


.PHONY: all clean build install install-deps install-deps-dev update-deps update-deps-dev test test-with-coverage test-with-coverage-json-output test-with-coverage-profile generate-test-coverage-html test-with-coverage-profile-html format-lint import-lint style-lint lint format-fix import-fix vet
