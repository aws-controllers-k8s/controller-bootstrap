# This configuration is a placeholder. Replace any values with relevant values for your
# service controller project.
---
annotations:
  capabilityLevel: Basic Install
  shortDescription: AWS {{ .ServiceID }} controller is a service controller for managing {{ .ServiceID }} resources
    in Kubernetes
displayName: AWS Controllers for Kubernetes - Amazon {{ .ServiceID }}
description: |-
  Manage Amazon {{ .ServiceID }} resources in AWS from within your Kubernetes cluster.


  **About Amazon {{ .ServiceID }}**


  {ADD YOUR DESCRIPTION HERE}


  **About the AWS Controllers for Kubernetes**


  This controller is a component of the [AWS Controller for Kubernetes](https://github.com/aws/aws-controllers-k8s)
  project. This project is currently in **developer preview**.
samples:
- kind: ExampleCustomKind
  spec: '{}'
- kind: SecondExampleCustomKind
  spec: '{}'
maintainers:
- name: "{{ .ServicePackageName }} maintainer team"
  email: "ack-maintainers@amazon.com"
links:
- name: Amazon {{ .ServiceID }} Developer Resources
  url: "{YOUR SERVICE DEVELOPER RESOURCES URL}"
