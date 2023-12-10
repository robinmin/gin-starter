GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
APP_NAME=gin-stater
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
	go install github.com/cosmtrek/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

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
	GO111MODULE=on $(GOCMD) build -o ./dist/bin/$(APP_NAME) ./cmd/server/main.go
#	GO111MODULE=on $(GOCMD) build -o ./dist/bin/$(APP_NAME)-client ./cmd/client/main.go
	@chmod u+x dist/bin/$(APP_NAME)
#	@chmod u+x dist/bin/$(APP_NAME)-client

#	errcheck -ignoretests ./cmd/client/main.go
#	go vet ./cmd/client/main.go
#	golangci-lint run -v ./cmd/client/main.go

# 	go build -ldflags="-s -w" -o cmd/my-app ./cmd/main.go

.PHONY: clean
clean: ## Remove all binaries.
	rm -rf ./dist/bin/
	rm -rf ./tmp/build-errors.log

.PHONY: test
test: ## Run the tests of the project.
	$(info ******************** running tests ********************)
	go test -coverprofile=log/coverage_${APP_VERSION}.out -v ./...
	go tool cover -html=log/coverage_${APP_VERSION}.html

# .PHONY: check
# check: ## Run precke before release
# 	$(info ******************** checking before commit ********************)
# 	goreleaser --snapshot --skip-publish --clean

.PHONY: lint
lint:  ## Run all available linters.
	$(info ******************** running lint tools ********************)
	errcheck -ignoretests ./cmd/server/main.go
	go vet ./cmd/server/main.go
	golangci-lint run -v ./cmd/server/main.go

# .PHONY: release
# release: ## Check before release
# 	$(info ******************** checking before release ********************)
# 	goreleaser release --clean

# .PHONY: generate
# generate: ## Generate code from SQL file by sqlc
# 	$(info ******************** generating code from sql ********************)
# 	sqlc generate
