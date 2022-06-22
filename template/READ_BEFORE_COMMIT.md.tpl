# Instructions

This document lays out the next steps for AWS service teams upon bootstrapping the ACK {{ .ServicePackageName }}-controller repository successfully.

1. To get started, look inside the `generator.yaml` file, and remove the resource(s) from the ignore list and execute `make build-controller` from the ACK code-generator.
2. Add any custom code inside the `templates` directory.
3. Refer to testing documentation to add e2e test inside the `test` directory.
4. Add your team members' GitHub aliases to your team controller alias in `OWNERS_ALIASES` file.
