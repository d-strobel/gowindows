# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Source .env file if available
ifneq ("$(wildcard .env)","")
	include .env
endif

# Format code
.PHONY: format
format:
	@printf "$(OK_COLOR)==> Format code$(NO_COLOR)\n"
	@go fmt ./...

# Download dependencies
.PHONY: dependencies
dependencies:
	@printf "$(OK_COLOR)==> Install dependencies$(NO_COLOR)\n"
	@go get -d -v ./...

# Setup requirements
.PHONY: vagrant-up
vagrant-up:
	@printf "$(OK_COLOR)==> Setup requirements$(NO_COLOR)\n"
	@$(MAKE) -C vagrant vagrant-up

# Remove requirements
.PHONY: vagrant-down
vagrant-down:
	@printf "$(OK_COLOR)==> Remove requirements$(NO_COLOR)\n"
	@$(MAKE) -C vagrant vagrant-down

# Unit tests
.PHONY: test
test: dependencies
	@printf "$(OK_COLOR)==> Run unit tests$(NO_COLOR)\n"
	@go test -short ./...

# Acceptance tests
.PHONY: testacc
testacc: dependencies
	@printf "$(OK_COLOR)==> Run acceptance tests$(NO_COLOR)\n"
	@go test ./...
