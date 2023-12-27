SHELL=/bin/bash -e -o pipefail

# constants
GOLANGCI_VERSION = 1.51.0
GOLANGCI_LINT = bin/golangci-lint-$(GOLANGCI_VERSION)

out:
	@mkdir -pv "$(@)"

download: ## Downloads the dependencies
	@go mod download

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b bin v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint "$(@)"

lint: fmt $(GOLANGCI_LINT) download ## Lints all code with golangci-lint
	@$(GOLANGCI_LINT) run

lint-reports: out/lint.xml

fmt: ## Formats all code with go fmt
	@go fmt ./...

test: ## Runs all tests
	@go test -p 1 -v ./...

run: export API_LOG_FORMATTER=text
run: export API_LOG_LEVEL=debug
run: ## Runs the application at 8080 port
	@go run cmd/api/main.go

coverage: out/report.json ## Displays coverage per func on cli
	@go tool cover -func=out/cover.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	@go tool cover -html=out/cover.out

test-reports: out/report.json

.PHONY: out/report.json
out/report.json: out
	go test ./... -coverprofile=out/cover.out --json | tee "$(@)"

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
