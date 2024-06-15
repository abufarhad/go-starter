PROJECT_NAME := go_starter
PKG_LIST := $(shell go list ${PROJECT_NAME}/tests/testing/... | grep -v /vendor/)


.PHONY: all dep build core test

all: build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


########################
### DEVELOP and TEST ###
########################
development:  ## Develop new build
	# booting up dependency containers
	@docker-compose up -d --build  db

	# building user_api
	@docker-compose up -d --build ${PROJECT_NAME}

test: ## Run unittests
	@go test -cover -short ${PKG_LIST} -v

coverage: ## Generate global code coverage report
	@go tool cover -func=cov.out

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)
	@docker-compose down

run: ## Run application
	@swag init
	@go run main.go serve