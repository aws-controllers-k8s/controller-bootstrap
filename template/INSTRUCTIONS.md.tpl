# Instructions

This document lays out the next steps for AWS service teams upon bootstrapping the ACK {{ .ServicePackageName }}-controller repository successfully.

1. To get started, look inside the `generator.yaml` file, ignore the resource(s) from the resource list and execute `make build-controller` from the ACK code-generator.
2. Add any custom code inside the `templates` directory, add e2e test inside the `test` directory.
3. Add your team members to your team controller alias in [OWNERS_ALIASES](OWNERS_ALIASES) file.
