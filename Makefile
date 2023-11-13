# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Source .env file if available
ifneq ("$(wildcard .env)","")
	include .env
endif

# Assert environment variables for testing
.PHONY: assert-test-environment
assert-test-environment:
	@printf "$(OK_COLOR)==> Assert test environment$(NO_COLOR)\n"

ifndef GOWINDOWS_TEST_SSH_HOST
	$(error GOWINDOWS_TEST_SSH_HOST is not set.)
endif
ifndef GOWINDOWS_TEST_SSH_PORT
	$(error GOWINDOWS_TEST_SSH_PORT is not set.)
endif
ifndef GOWINDOWS_TEST_SSH_USERNAME
	$(error GOWINDOWS_TEST_SSH_USERNAME is not set.)
endif
ifndef GOWINDOWS_TEST_SSH_PASSWORD
	$(error GOWINDOWS_TEST_SSH_PASSWORD is not set.)
endif

# Format code
.PHONY: format
format:
	@printf "$(OK_COLOR)==> Format code$(NO_COLOR)\n"
	go fmt ./...

# Download dependencies
.PHONY: dependencies
dependencies:
	@printf "$(OK_COLOR)==> Install dependencies$(NO_COLOR)\n"
	go get -d -v ./...

# Download requirements
.PHONY: requirements
requirements:
	@printf "$(OK_COLOR)==> Install requirements$(NO_COLOR)\n"

# Unit tests
.PHONY: test
test: dependencies
	@printf "$(OK_COLOR)==> Run unit tests$(NO_COLOR)\n"
	go test ./...

# Acceptance tests
.PHONY: testacc
testacc: assert-test-environment requirements dependencies
	@printf "$(OK_COLOR)==> Run acceptance tests$(NO_COLOR)\n"
	go test ./...
