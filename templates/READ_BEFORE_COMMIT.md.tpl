# Instructions

This document lays out the next steps for AWS service teams upon bootstrapping the ACK {{ .ServicePackageName }}-controller repository successfully.

1. To get started, edit `generator.yaml` file and comment out (with a `#`) each of the resource(s) from the ignore list, generating your controller after every edit.
- This step will enable each resource within the service to be generated into an ACK custom resource definition. As you generate each of them, the code-generator may require additional configuration (in the `generator.yaml` file) in order for you to continue.

2. Update fields under `service` inside `metadata.yaml`.
- This file provides display information about the controller and the service it supports. This file is intended to be manually updated to ensure it matches the proper casing and terminology for the service.
For `full_name`, provide the full display name of the service (eg. `Amazon Simple Storage Service`). For `short_name`, provide the abbreviation or shortened service name (eg. `S3` or `SageMaker` - do not include `AWS` or `Amazon`).
For `link`, provide a link to the homepage for the service. For `documentation`, provide a link to the homepage for the service documentation/user guide.

3. Add your team members' GitHub aliases to the `OWNERS_ALIASES` file.
- This file determines who has permissions to review and merge code within the repository.

4. Add any custom code to the `templates` directory.
- As you work on the controller, you may need to add custom hooks to support functionality specific to your service. Hook code for any resource should be placed in a file with the path `templates/hooks/<resource name>/<hook name>.go.tpl`.

5. Refer to the open source documentation on testing for creating and running tests, which should be added in the `test/e2e` directory.
- All PRs containing changes to CRDs should be accompanied with appropriate tests for any new functionality or custom code. Documentation for these tests can be found in the "Contributor Docs" section of the ACK website.
