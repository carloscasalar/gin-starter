SHELL=/bin/bash -e -o pipefail

# constants
REVIVE_VERSION = v1.3.4
DOCKER_REPO = api
DOCKER_TAG = latest

out:
	@mkdir -pv "$(@)"

download: ## Downloads the dependencies
	@go mod download

build: out/bin ## Builds all binaries

GO_BUILD = mkdir -pv "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...
.PHONY: out/bin
out/bin:
	$(GO_BUILD)

docker: ## Builds docker image
	docker buildx build -t $(DOCKER_REPO):$(DOCKER_TAG) .

lint: ## Lints all code with revive
	@go install github.com/mgechev/revive@$(REVIVE_VERSION)
	@revive -config revive.toml -formatter friendly ./...

lint-reports: out/lint.xml

fmt: ## Formats all code with go fmt
	@go fmt ./...

test: ## Runs all tests
	@go test -p 1 -v ./...

run: export API_LOG_FORMATTER=text
run: export API_LOG_LEVEL=debug
run: ## Runs the application at 8080 port
	@go run cmd/api/main.go

run-docker: fmt build docker ## Runs the trick inside docker
	docker run --rm -it -p8080:8080 --env API_LOG_FORMATTER="text" $(DOCKER_REPO):$(DOCKER_TAG)

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
