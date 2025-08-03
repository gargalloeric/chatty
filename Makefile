# ============================================================================== #
# HELPERS
# ============================================================================== #

## help: prints this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans  && [ $${ans:-N} = y ] 

# ============================================================================== #
# DEVELOPMENT
# ============================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api

## run/client: run the cmd/client application
.PHONY: run/client
run/client:
	@go run ./cmd/client
