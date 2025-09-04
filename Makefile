.DEFAULT_GOAL := help

PROJECT_ROOT:=$(shell git rev-parse --show-toplevel)

# Load env properties , db name, port, etc...
# nb: You can change the default config with `make ENV_CONTEXT=".env.uat" `
ENV_CONTEXT ?= .env.local
ENV_CONTEXT_PATH:=$(PROJECT_ROOT)/$(ENV_CONTEXT)

## Override any default values in the parent .env, with your own
-include $(ENV_CONTEXT_PATH)

COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
REPO := $(shell basename `git rev-parse --show-toplevel`)
DATE := $(shell date +%Y-%m-%d-%H-%M-%S)
APP_NAME := $(REPO)

GO_BUILD_FLAGS=-ldflags=""

BUILD_PATH ?= "cmd"
BINARY_PATH ?= "dist"

ci: log scan test ## Run CI checks

test: ## Run unit tests
	go test --short -cover -failfast ./...

test-watch: ## Run unit tests in watch mode
	gow test -v --short -cover -failfast ./...

scan: ## run security scan
	govulncheck ./...
	golangci-lint run ./...

tools-get: ## Get project tools required
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)



#####################
#  Private Targets  #
#####################

log: # log env vars
	@echo "\n"
	@echo "COMMIT               $(COMMIT)"
	@echo "BRANCH               $(BRANCH)"
	@echo "APP_NAME             $(APP_NAME)"
	@echo "REPO                 $(REPO)"
	@echo "DATE                 $(DATE)"
	@echo "ENVIRONMENT          $(ENVIRONMENT)"
	@echo "LOG_LEVEL            $(LOG_LEVEL)"
	@echo "\n"
