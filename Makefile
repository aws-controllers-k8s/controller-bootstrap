SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

GO_CMD_FLAGS=-tags codegen
AWS_SERVICE=$(shell echo $(SERVICE))
SERVICE_MODEL_NAME=$(shell echo $(MODEL_NAME))
ifeq ($(SERVICE_MODEL_NAME),)
  SERVICE_MODEL_NAME:=""
endif
ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CONTROLLER_BOOTSTRAP=./bin/controller-bootstrap
CONTROLLER_DIR=${ROOT_DIR}/../${AWS_SERVICE}-controller
CODE_GEN_DIR=${ROOT_DIR}/../code-generator
AWS_SDK_GO_VERSION=$(shell curl -H "Accept: application/vnd.github.v3+json" \
                                 https://api.github.com/repos/aws/aws-sdk-go/releases/latest | jq -r '.tag_name')
ACK_RUNTIME_VERSION=$(shell curl -H "Accept: application/vnd.github.v3+json" \
								 https://api.github.com/repos/aws-controllers-k8s/runtime/releases/latest | jq -r '.tag_name')
.DEFAULT_GOAL=run
EXISTING_CONTROLLER="true"
DRY_RUN="false"

.PHONY: build, generate, update, init, run, clean

build:
	@go build ${GO_CMD_FLAGS} -o ${CONTROLLER_BOOTSTRAP} ./cmd/controller-bootstrap/*.go

generate: build
	@${CONTROLLER_BOOTSTRAP} generate -s ${AWS_SERVICE} -r ${ACK_RUNTIME_VERSION} -v ${AWS_SDK_GO_VERSION} -d=${DRY_RUN} -o ${ROOT_DIR}/../${AWS_SERVICE}-controller -m ${SERVICE_MODEL_NAME}

update: build
	@${CONTROLLER_BOOTSTRAP} generate -s ${AWS_SERVICE} -r ${ACK_RUNTIME_VERSION} -v ${AWS_SDK_GO_VERSION} -d=${DRY_RUN} -e=${EXISTING_CONTROLLER} -o ${ROOT_DIR}/../${AWS_SERVICE}-controller -m ${SERVICE_MODEL_NAME}

init: generate
	@export SERVICE=${AWS_SERVICE}
	@echo "build controller attempt #1"
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@echo "missing go.sum entry, running go mod tidy"
	@cd ${CONTROLLER_DIR} && go mod tidy
	@echo "build controller attempt #2"
	@cd ${CODE_GEN_DIR} && make -i build-controller >/dev/null 2>/dev/null
	@echo "go.sum outdated, running go mod tidy"
	@cd ${CONTROLLER_DIR} && go mod tidy
	@echo "final build controller attempt"
	@cd ${CODE_GEN_DIR} && make build-controller
	@echo "look inside ${SERVICE}-controller/INSTRUCTIONS.md for further instructions"

run:
	@if [ -f ${CONTROLLER_DIR}/cmd/controller/main.go ]; then \
	  make update; \
	else \
	  make init; \
	fi

clean:
	@cd ${CONTROLLER_DIR}
	@rm -rf ${CONTROLLER_DIR}/..?* ${CONTROLLER_DIR}/.[!.]* ${CONTROLLER_DIR}/*
