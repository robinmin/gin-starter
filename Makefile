GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
APP_NAME=gin-starter
APP_VERSION?=0.0.1
SERVICE_PORT?=3000

FILES=$(shell find . -name "*.go")

###############################################################################
.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install
install: ## Install dependencies.
	$(info ******************** installing dependencies ********************)
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install github.com/pressly/goose/v3/cmd/goose@latest
	$(GOCMD) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: start
start: ## Start dev server.
	$(info ******************** running dev server ********************)
	air

.PHONY: all
all: fmt build test lint

.PHONY: fmt
fmt: ## Format all code.
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(FILES)) || (gofmt -d $(FILES); exit 1)

.PHONY: build
build: ## Build binaries
	$(info ******************** building binaries ********************)
	GO111MODULE=on $(GOCMD) build -ldflags="-s -w" -o ./dist/bin/$(APP_NAME) ./cmd/cli/main.go
	@chmod u+x dist/bin/$(APP_NAME)

.PHONY: clean
clean: ## Remove all binaries.
	rm -rf ./dist/bin/
	rm -rf ./tmp/build-errors.log

.PHONY: test
test: ## Run the tests of the project.
	$(info ******************** running tests ********************)
	$(GOCMD) test -coverprofile=log/coverage_${APP_VERSION}.out -v ./...
	$(GOCMD) tool cover -html=log/coverage_${APP_VERSION}.html

.PHONY: check
check: ## Run precke before committing
	$(info ******************** checking before committing ********************)
# 	use git commit hook to check before commiting the code into git
	@.git/hooks/prepare-commit-msg .

#	goreleaser --snapshot --skip-publish --clean

.PHONY: lint
lint:  ## Run all available linters.
	$(info ******************** running lint tools ********************)
	errcheck -ignoretests ./cmd/cli/main.go
	$(GOCMD) vet ./cmd/cli/main.go
	golangci-lint run -v ./cmd/cli/main.go

# .PHONY: release
# release: ## Check before release
# 	$(info ******************** checking before release ********************)
# 	goreleaser release --clean

.PHONY: generate
generate: ## Generate code from SQL file by sqlc
	$(info ******************** generating code from sql ********************)
	sqlc -f ./schema/.sqlc.yaml generate
