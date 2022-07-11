#!/usr/bin/env bash

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
CONTROLLER_BOOTSTRAP_DIR="$SCRIPTS_DIR/.."
CONTROLLER_BOOTSTRAP="$CONTROLLER_BOOTSTRAP_DIR/bin/controller-bootstrap"
CONTROLLER_BOOTSTRAP_TEMPLATES_DIR="$CONTROLLER_BOOTSTRAP_DIR/templates"
SERVICE="eks"
CONTROLLER_NAME="$SERVICE-controller"
CONTROLLER_DIR="$CONTROLLER_BOOTSTRAP_DIR/test_output/$CONTROLLER_NAME"

export SERVICE
export CONTROLLER_DIR
cd $CONTROLLER_BOOTSTRAP_DIR
make generate

if [[ ! -e "$CONTROLLER_DIR/TEST_FILE" ]]; then
  echo "test_generate.sh][INFO] 'TEST_FILE' is not found in $CONTROLLER_NAME. Creating a new template test file, controller-bootstrap/templates/TEST_FILE.tpl ..."
  cd "$CONTROLLER_BOOTSTRAP_TEMPLATES_DIR"
  touch "TEST_FILE.tpl"
  rm -rf ${CONTROLLER_DIR}
  echo "Generating $CONTROLLER_NAME using command 'make generate'"
  cd "$CONTROLLER_BOOTSTRAP_DIR"
  make generate
fi

if [[ ! -e "$CONTROLLER_DIR/TEST_FILE" ]]; then
  echo "test_generate.sh][ERROR] Unable to find 'TEST_FILE' in $CONTROLLER_NAME. Failed to generate the controller. Exiting"
  exit 1
else
  echo "Generated successfully and found 'TEST_FILE' in $CONTROLLER_NAME"
fi

cd "$CONTROLLER_BOOTSTRAP_TEMPLATES_DIR"
rm "TEST_FILE.tpl"
