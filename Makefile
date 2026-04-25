unittest-json:
	@mkdir -p build
	@go test -v -json -coverprofile=build/coverage.out -covermode=atomic ./... -timeout 20s -failfast > build/unit-test-output.json

coverage-report:
	@go tool cover -html=build/coverage.out -o build/coverage.html
	@echo "### Code Coverage: $$(go tool cover -func=build/coverage.out | awk '/^total:/{print $$3}')"

coverage:
	@go tool cover -html=build/coverage.out -o build/coverage.html
	@go tool cover -func=build/coverage.out | awk 'END {print $$3}'

vulncheck-sarif:
	@mkdir -p build
	@go tool -modfile=./tools/go.mod govulncheck -format sarif ./... > build/govulncheck-report.sarif

build:
	@go build ./...
