# Source .env file if available
ifneq ("$(wildcard .env)","")
	include .env
endif

# Ensure prerequesites
testacc: assert-test-environment

# Assert environment variables for testing
.PHONY: assert-test-environment
assert-test-environment:
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

# Acceptance tests
.PHONY: testacc
testacc:
	go test ./...
