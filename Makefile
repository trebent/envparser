GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin

build:
@go build ./...

coverage:
@go tool cover -html=build/coverage.out -o build/coverage.html
@go tool cover -func=build/coverage.out | awk 'END {print $$3}'

coverage/report:
@go tool cover -html=build/coverage.out -o build/coverage.html
@echo "### Code Coverage: $$(go tool cover -func=build/coverage.out | awk '/^total:/{print $$3}')"

install/lint:
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOBIN) v2.12.2

static-analysis/lint:
golangci-lint run --fix

static-analysis/vulncheck:
@go tool -modfile=./tools/go.mod govulncheck ./...

static-analysis/vulncheck-sarif:
@mkdir -p build
@go tool -modfile=./tools/go.mod govulncheck -format sarif ./... > build/govulncheck-report.sarif

test/unit:
@mkdir -p build
@go test -v ./... -timeout 20s -failfast -coverprofile=build/coverage.out -covermode=atomic

test/unit-json:
@mkdir -p build
@go test -v -json -coverprofile=build/coverage.out -covermode=atomic ./... -timeout 20s -failfast > build/unit-test-output.json
