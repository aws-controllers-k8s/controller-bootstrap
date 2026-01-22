SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

# Allow GITHUB_TOKEN to be passed in via CLI or environment
# Usage: make generate GITHUB_TOKEN=ghp_yourtokenhere
GITHUB_TOKEN ?=

# Helper to construct the Authorization header if the token exists
AUTH_HEADER := $(if $(GITHUB_TOKEN),-H "Authorization: Bearer $(GITHUB_TOKEN)",)
GITHUB_API_H := -H "Accept: application/vnd.github.v3+json" $(AUTH_HEADER)

AWS_SERVICE=$(shell echo $(SERVICE))
SERVICE_MODEL_NAME=$(shell echo $(MODEL_NAME))
ifeq ($(SERVICE_MODEL_NAME),)
  SERVICE_MODEL_NAME=""
endif

ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CONTROLLER_BOOTSTRAP=./bin/controller-bootstrap
CODE_GEN_DIR=${ROOT_DIR}/../code-generator

CONTROLLER_DIR:=$(or $(CONTROLLER_DIR),${ROOT_DIR}/../${AWS_SERVICE}-controller)
ACK_RUNTIME_VERSION:=$(or $(ACK_RUNTIME_VERSION),$(shell curl -s $(GITHUB_API_H) \
    https://api.github.com/repos/aws-controllers-k8s/runtime/releases/latest | jq -r '.tag_name'))
AWS_SDK_GO_VERSION:=$(or $(AWS_SDK_GO_VERSION),$(shell curl -s $(GITHUB_API_H) \
    https://api.github.com/repos/aws/aws-sdk-go/releases/latest | jq -r '.tag_name'))
TEST_INFRA_COMMIT_SHA:=$(or $(TEST_INFRA_COMMIT_SHA),$(shell curl -s $(GITHUB_API_H) \
    https://api.github.com/repos/aws-controllers-k8s/test-infra/commits | jq -r ".[0].sha"))

.DEFAULT_GOAL=run
DRY_RUN="false"
REFRESH_CACHE="true"

# Build ldflags
VERSION ?= "v0.0.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS=-ldflags "-X main.version=$(VERSION) \
			-X main.buildHash=$(GITCOMMIT) \
			-X main.buildDate=$(BUILDDATE)"

# We need to use the codegen tag when building and testing because the
# aws-sdk-go/private/model/api package is gated behind a build tag "codegen"...
GO_CMD_FLAGS=-tags codegen

.PHONY: build, generate, init, run, clean

build:
	@go build ${GO_CMD_FLAGS} -o ${CONTROLLER_BOOTSTRAP} ./cmd/controller-bootstrap/main.go

generate: build
# 	log variables for debugging
	@echo "Generating ${AWS_SERVICE}-controller with the following parameters:"
	@echo "  CONTROLLER_DIR=${CONTROLLER_DIR}"
	@echo "  ACK_RUNTIME_VERSION=${ACK_RUNTIME_VERSION}"
	@echo "  AWS_SDK_GO_VERSION=${AWS_SDK_GO_VERSION}"
	@echo "  SERVICE_MODEL_NAME=${SERVICE_MODEL_NAME}"	
	@${CONTROLLER_BOOTSTRAP} generate --aws-service-alias ${AWS_SERVICE} --ack-runtime-version ${ACK_RUNTIME_VERSION} \
    --aws-sdk-go-version ${AWS_SDK_GO_VERSION} --dry-run=${DRY_RUN} --output-path ${CONTROLLER_DIR} \
    --model-name ${SERVICE_MODEL_NAME} --refresh-cache=${REFRESH_CACHE} --test-infra-commit-sha ${TEST_INFRA_COMMIT_SHA}

init: generate
	@export SERVICE=${AWS_SERVICE}
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@cd ${CONTROLLER_DIR} && go mod tidy
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@cd ${CONTROLLER_DIR} && go mod tidy
	@cd ${CODE_GEN_DIR} && make build-controller
	@echo "${AWS_SERVICE}-controller generated successfully, look inside ${AWS_SERVICE}-controller/READ_BEFORE_COMMIT.md for further instructions"

run:
	@if [ -f ${CONTROLLER_DIR}/cmd/controller/main.go ]; then \
        make build; \
        ${CONTROLLER_BOOTSTRAP} update --aws-service-alias ${AWS_SERVICE} --output-path ${CONTROLLER_DIR}; \
        echo "${AWS_SERVICE}-controller updated successfully with the latest templates of project description files"; \
	else \
	    make init; \
	fi

test: 				## Run code tests
	go test ${GO_CMD_FLAGS} ./...

clean:
	@rm -rf ${CONTROLLER_DIR}/..?* ${CONTROLLER_DIR}/.[!.]* ${CONTROLLER_DIR}/*
