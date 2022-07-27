# Controller-Bootstrap for AWS Controllers for Kubernetes (ACK)

The ACK controller-bootstrap tool automates the repository bootstrap process for new ACK service controllers in the `aws-controller-k8s` Github organization.
For existing ACK service controllers, this tool will update the service controllers with the latest templates of the project description files, 
which include: `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `GOVERNANCE.md`, `LICENSE`, `NOTICE`, `SECURITY.md` files.

## Getting Started

First, clone the `aws-controllers-k8s/controller-bootstrap` repository and run the following script in the controller-bootstrap.
```
export SERVICE=${AWS_SERVICE_NAME}
make
```
For a new ACK service controller, the `make` command bootstraps an ACK service controller repository. For an existing ACK service controller, it updates the service controller repository with the latest templates of the project description files.


To generate the common directories and files for a new ACK service controller using CLI command, the user can run the `generate` command from the controller-bootstrap.
```
controller-bootstrap generate --aws-service-alias ${AWS_SERVICE} --ack-runtime-version ${ACK_RUNTIME_VERSION}
    --aws-sdk-go-version ${AWS_SDK_GO_VERSION} --dry-run=${DRY_RUN} --output-path ${CONTROLLER_DIR}
    --model-name ${SERVICE_MODEL_NAME} --refresh-cache=${REFRESH_CACHE} --test-infra-commit-sha ${TEST_INFRA_COMMIT_SHA}
```

To update an existing ACK service controller with the latest templates of the project description files using CLI command, the user can run the `update` command from the controller-bootstrap.
```
controller-bootstrap update --aws-service-alias ${AWS_SERVICE} --output-path ${CONTROLLER_DIR}
```

The command-line arguments of the controller-bootstrap `generate` and `update` commands are described in the [Usage](#usage).

## Usage
```
Usage:
controller-bootstrap generate [flags]

Flags:
--ack-runtime-version string     Version of aws-controllers-k8s/runtime
--aws-sdk-go-version string      Version of github.com/aws/aws-sdk-go used to infer service metadata and resources
-h, --help                       help for generate
--model-name string              Optional: service model name of the corresponding service alias
--refresh-cache                  Optional: if true, and aws-sdk-go repo is already cloned, will git pull the latest aws-sdk-go commit (default true)
--test-infra-commit-sha string   Commit SHA of aws-controllers-k8s/test-infra

Global Flags:
--aws-service-alias string   AWS service alias
--dry-run                    Optional: if true, output files to stdout (default true)
--output-path string         Path to ACK service controller directory to bootstrap
```
```
Usage:
  controller-bootstrap update [flags]

Flags:
  -h, --help   help for update

Global Flags:
      --aws-service-alias string   AWS service alias
      --dry-run                    Optional: if true, output files to stdout (default true)
      --output-path string         Path to ACK service controller directory to bootstrap
```

## Community, discussion, contribution, and support

We welcome community contributions and pull requests.

See our [contribution guide](/CONTRIBUTING.md) for more information on how to
report issues, set up a development environment, and submit code.

### Code of conduct

Participation in the AWS community is governed by the [Amazon Open Source Code of Conduct](CODE_OF_CONDUCT.md).
