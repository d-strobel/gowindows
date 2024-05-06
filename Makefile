# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Source .env file if available
ifneq ("$(wildcard .env)","")
	include .env
else
	include ./vagrant/vagrant.env
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

.PHONY: generate
generate:
	@printf "$(OK_COLOR)==> Go generate$(NO_COLOR)\n"
	@go generate

# Setup requirements
.PHONY: vagrant-up
vagrant-up:
	@printf "$(OK_COLOR)==> Setup vagrant machines$(NO_COLOR)\n"
	@$(MAKE) -C vagrant vagrant-up

# Remove requirements
.PHONY: vagrant-down
vagrant-down:
	@printf "$(OK_COLOR)==> Remove vagrant machines$(NO_COLOR)\n"
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

.PHONY: check-env
check-env:
	@printf "$(OK_COLOR)==> Environment variables for default Windows test machine$(NO_COLOR)\n"
	@echo "Host: $(GOWINDOWS_TEST_HOST)"
	@echo "Username: $(GOWINDOWS_TEST_USERNAME)"
	@echo "Password: $(GOWINDOWS_TEST_PASSWORD)"
	@echo "SSH Port: $(GOWINDOWS_TEST_SSH_PORT)"
	@echo "SSH private key path to ed25519: $(GOWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH)"
	@echo "SSH private key path to rsa: $(GOWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH)"
	@echo "WinRM http port: $(GOWINDOWS_TEST_WINRM_HTTP_PORT)"
	@echo "WinRM https port: $(GOWINDOWS_TEST_WINRM_HTTPS_PORT)"
	@printf "\n$(OK_COLOR)==> Environment variables for Active-Directory Windows test machine$(NO_COLOR)\n"
	@echo "Host: $(GOWINDOWS_TEST_AD_HOST)"
	@echo "Username: $(GOWINDOWS_TEST_AD_USERNAME)"
	@echo "Password: $(GOWINDOWS_TEST_AD_PASSWORD)"
	@echo "SSH Port: $(GOWINDOWS_TEST_AD_SSH_PORT)"
	@echo "WinRM http port: $(GOWINDOWS_TEST_AD_WINRM_HTTP_PORT)"
	@echo "WinRM https port: $(GOWINDOWS_TEST_AD_WINRM_HTTPS_PORT)"
