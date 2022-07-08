#!/usr/bin/env bash

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
CONTROLLER_BOOTSTRAP_DIR="$SCRIPTS_DIR/.."
CONTROLLER_BOOTSTRAP="$CONTROLLER_BOOTSTRAP_DIR/bin/controller-bootstrap"
NOTICE_TPL_FILE="$CONTROLLER_BOOTSTRAP_DIR/templates/NOTICE.tpl"
AWS_SERVICE="eks"
CONTROLLER_NAME="$AWS_SERVICE-controller"
OUTPUT_DIR="$CONTROLLER_BOOTSTRAP_DIR/test_output"
CONTROLLER_DIR="$OUTPUT_DIR/$CONTROLLER_NAME"

if ! grep -wq -- "this line is added by ACK update test" "$CONTROLLER_DIR/NOTICE"; then
    echo "Unable to find 'this line is added by ACK update test' in the 'NOTICE' file of $CONTROLLER_NAME. Adding 'this line is added by ACK update test' in controller-bootstrap/templates/NOTICE.tpl file ..."
    echo "this line is added by ACK update test" >> $NOTICE_TPL_FILE
    echo "Updating the project description files in the existing $CONTROLLER_NAME"
    cd "$CONTROLLER_BOOTSTRAP_DIR"
    make build
    ${CONTROLLER_BOOTSTRAP} update --aws-service-alias ${AWS_SERVICE} --output-path ${CONTROLLER_DIR}
fi

if ! grep -wq -- "this line is added by ACK update test" "$CONTROLLER_DIR/NOTICE"; then
    echo "Unable to find 'this line is added by ACK update test' in the 'NOTICE' file. Failed to update the project description files in $CONTROLLER_NAME. Exiting"
    rm -rf ${OUTPUT_DIR}
    exit 1
else
    echo "Updated successfully and found 'this line is added by ACK update test' in the 'NOTICE' file of $CONTROLLER_NAME"
fi

sed -i '' "/this line is added by ACK update test/d" $NOTICE_TPL_FILE
rm -rf ${OUTPUT_DIR}
