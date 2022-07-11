#!/usr/bin/env bash

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
CONTROLLER_BOOTSTRAP_DIR="$SCRIPTS_DIR/.."
CONTROLLER_BOOTSTRAP="$CONTROLLER_BOOTSTRAP_DIR/bin/controller-bootstrap"
NOTICE_TPL_FILE="$CONTROLLER_BOOTSTRAP_DIR/templates/NOTICE.tpl"
SERVICE="eks"
CONTROLLER_NAME="$SERVICE-controller"
OUTPUT_DIR="$CONTROLLER_BOOTSTRAP_DIR/test_output"
CONTROLLER_DIR="$OUTPUT_DIR/$CONTROLLER_NAME"
TEXT_TO_FIND="'this line is added by ACK update test'"

if ! grep -wq -- "$TEXT_TO_FIND" "$CONTROLLER_DIR/NOTICE"; then
    echo "test_update.sh][DEBUG] Unable to find $TEXT_TO_FIND in the 'NOTICE' file of $CONTROLLER_NAME. Adding $TEXT_TO_FIND in controller-bootstrap/templates/NOTICE.tpl file ..."
    echo "$TEXT_TO_FIND" >> $NOTICE_TPL_FILE
    echo "Updating the project description files in the existing $CONTROLLER_NAME"
    cd "$CONTROLLER_BOOTSTRAP_DIR"
    make build
    ${CONTROLLER_BOOTSTRAP} update --aws-service-alias ${SERVICE} --output-path ${CONTROLLER_DIR}
fi

if ! grep -wq -- "$TEXT_TO_FIND" "$CONTROLLER_DIR/NOTICE"; then
    echo "test_update.sh][ERROR] Unable to find "$TEXT_TO_FIND" in the 'NOTICE' file. Failed to update the project description files in $CONTROLLER_NAME. Exiting"
    rm -rf ${OUTPUT_DIR}
    exit 1
else
    echo "Updated successfully and found "$TEXT_TO_FIND" in the 'NOTICE' file of $CONTROLLER_NAME"
fi

sed -i '' "/$TEXT_TO_FIND/d" $NOTICE_TPL_FILE
rm -rf ${OUTPUT_DIR}
