GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
APP_NAME=gin-starter
APP_VERSION?=0.0.1
# SERVICE_PORT?=3000

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
MAKE_COLOR=\033[33;01m%-20s\033[0m

# This repo's root import path (under GOPATH).
PKG := github.com/robinmin/gin-starter

GO_FILES=$(shell find . -name "*.go" | grep -v vendor | uniq)
BUILD_DATE = $(shell date -u '+%Y.%m.%d')
GOLANG_VERSION ?= $(shell go version | cut -d" " -f3 | sed 's/go//')
GIT_REV    ?= $(shell git rev-parse --short HEAD)
GIT_TAG    ?= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_BRANCH ?= $(shell git branch|grep '*'| cut -f2 -d' ')
GIT_DIRTY  ?= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

###############################################################################
.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: version
version: ## show build info
	@echo "$(OK_COLOR) ******************** Build info ******************** $(NO_COLOR)"
	@echo "Application Name:    ${APP_NAME}"
	@echo "Application Version: ${APP_VERSION}"
	@echo "Golang Version:      ${GOLANG_VERSION}"
	@echo "Date:                ${BUILD_DATE}"
	@echo "Git Tag:             ${GIT_TAG}"
	@echo "Git Rev:             ${GIT_REV}"
	@echo "Git Tree State:      ${GIT_DIRTY}"


.PHONY: install
install: ## Install dependencies.
	@echo "$(OK_COLOR) ******************** Installing dependencies ******************** $(NO_COLOR)"
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install github.com/pressly/goose/v3/cmd/goose@latest
	$(GOCMD) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	pre-commit install --hook-type commit-msg

.PHONY: install-tools
install-tools: ## install dev tools, linters, code generaters, etc
	@echo "$(OK_COLOR) ******************** Installing tools from tools/tools.go ******************** $(NO_COLOR)"
	@export GOBIN=$$PWD/tools/bin; export PATH=$$GOBIN:$$PATH; cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: all
# all: fmt build test lint
all:  install-tools lint build test ## build, lint and test

.PHONY: ci-all
ci-all: install-tools lint build test-ci ## run build, lint and test pipeline

.PHONY: start
start: ## Start dev server.
	@echo "$(OK_COLOR) ******************** Running dev server ******************** $(NO_COLOR)"
	air

.PHONY: run
run: ## run application locally with the given .env file
	@echo "$(OK_COLOR) ******************** Running application ******************** $(NO_COLOR)"
	@(sh -ac 'source .env && go run cmd/cli/main.go')

.PHONY: fmt
fmt: ## Format all code.
	@echo "$(OK_COLOR) ******************** Formatting go files ********************$(NO_COLOR)"
	@tools/bin/gofumpt -l -w $(GO_FILES)
#	@tools/bin/gci write $(GO_FILES) -s "standard, default, Prefix($(PKG))"
	@tools/bin/gci write $(GO_FILES)

.PHONY: build
build: ## Build binaries
	@echo "$(OK_COLOR) ******************** Building binaries ******************** $(NO_COLOR)"
	GO111MODULE=on $(GOCMD) build -ldflags="-s -w" -o ./dist/bin/$(APP_NAME) ./cmd/cli/main.go
	@chmod u+x dist/bin/$(APP_NAME)

.PHONY: build-docker
build-docker: ## build application staically for docker
	@echo "$(OK_COLOR) ******************** Building binaries for docker ******************** $(NO_COLOR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -ldflags="-s -w" -a -o ./dist/bin/$(APP_NAME)  ./cmd/cli/main.go

.PHONY: clean
clean: ## Remove all binaries.
	@echo "$(OK_COLOR) ******************** Cleanup binaries ******************** $(NO_COLOR)"
	rm -rf ./dist/bin/
	rm -rf ./tools/bin/

.PHONY: test
# test: test-unit test-integration test-e2e ## run all tests
test: test-unit ## run all tests

# .PHONY: test-docker
# test-docker: ## run all tests in with docker-compose
# 	docker-compose up --build --force-recreate --no-deps -d db
# 	docker-compose run migrate
# 	@(sh -ac 'source .env && make test')
# 	docker-compose down --remove-orphans --volumes

.PHONY: test-unit
test-unit: ## run unit tests
	@echo "$(OK_COLOR) ******************** Running unit tests ********************$(NO_COLOR)"
	go test --race --count=1 ./...

# .PHONY: test-integration
# test-integration: ## run integration tests
# 	@echo "$(OK_COLOR) ******************** Running integration tests ******************** $(NO_COLOR)"
# 	go test --tags "integration" --race --count=1 ./tests/integration/...

.PHONY: test-e2e
test-e2e: ## run e2e tests
	@echo "$(OK_COLOR) ******************** Running E2E tests ******************** $(NO_COLOR)"
	go test --tags "acceptance" --race --count=1 ./tests/e2e/...

.PHONY: test-ci
test-ci: ## runing all tests with coverage
	@echo "$(OK_COLOR) ******************** Generating code coverage ******************** $(NO_COLOR)"
	sh tools/generate-fake-tests.sh
	sh tools/coverage.sh

.PHONY: check
check: ## Run precke before committing
	@echo "$(OK_COLOR) ******************** Checking before committing ******************** $(NO_COLOR)"
	pre-commit run --all-files
#	goreleaser --snapshot --skip-publish --clean

.PHONY: lint
lint:  ## Run all available linters.
	@echo "$(OK_COLOR) ******************** Running lint tools ******************** $(NO_COLOR)"
	errcheck -ignoretests ./cmd/cli/main.go
	$(GOCMD) vet ./cmd/cli/main.go
	golangci-lint run -v ./cmd/cli/main.go

# .PHONY: release
# release: ## Check before release
# 	@echo "$(OK_COLOR) ******************** checking before release ******************** $(NO_COLOR)"
# 	goreleaser release --clean

.PHONY: generate
generate: ## Generate code from SQL file by sqlc
	@echo "$(OK_COLOR) ******************** Generating code from sql ******************** $(NO_COLOR)"
	sqlc -f ./schema/.sqlc.yaml generate
