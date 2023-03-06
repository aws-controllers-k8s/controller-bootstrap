module github.com/aws-controllers-k8s/{{ .Service.Name.Lower }}-controller

go 1.19

require (
	github.com/aws-controllers-k8s/runtime v0.24.0
	github.com/aws/aws-sdk-go v1.44.214
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.11.0
)
