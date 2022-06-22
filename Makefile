SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

AWS_SERVICE=$(shell echo $(SERVICE))
SERVICE_MODEL_NAME=$(shell echo $(MODEL_NAME))
ifeq ($(SERVICE_MODEL_NAME),)
  SERVICE_MODEL_NAME=""
endif

ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CONTROLLER_DIR=${ROOT_DIR}/../${AWS_SERVICE}-controller
CONTROLLER_BOOTSTRAP=./bin/controller-bootstrap
CODE_GEN_DIR=${ROOT_DIR}/../code-generator

AWS_SDK_GO_VERSION=$(shell curl -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/aws/aws-sdk-go/releases/latest | jq -r '.tag_name')
ACK_RUNTIME_VERSION=$(shell curl -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/aws-controllers-k8s/runtime/releases/latest | jq -r '.tag_name')

.DEFAULT_GOAL=run
DRY_RUN="false"
EXISTING_CONTROLLER="false"

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
	@${CONTROLLER_BOOTSTRAP} generate -s ${AWS_SERVICE} -r ${ACK_RUNTIME_VERSION} -v ${AWS_SDK_GO_VERSION} -d=${DRY_RUN} -e=${EXISTING_CONTROLLER} -o ${ROOT_DIR}/../${AWS_SERVICE}-controller -m ${SERVICE_MODEL_NAME}

init: generate
	@export SERVICE=${AWS_SERVICE}
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@cd ${CONTROLLER_DIR} && go mod tidy
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@cd ${CONTROLLER_DIR} && go mod tidy
	@cd ${CODE_GEN_DIR} && make build-controller
	@echo "${AWS_SERVICE}-controller generated successfully, look inside ${AWS_SERVICE}-controller/INSTRUCTIONS.md for further instructions"

run:
	@if [ -f ${CONTROLLER_DIR}/cmd/controller/main.go ]; then \
  	  	EXISTING_CONTROLLER="true"; \
	    make generate; \
	else \
	    make init; \
	fi

clean:
	@rm -rf ${CONTROLLER_DIR}/..?* ${CONTROLLER_DIR}/.[!.]* ${CONTROLLER_DIR}/*
