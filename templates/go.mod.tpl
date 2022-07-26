module github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller

go 1.17

require (
	github.com/aws-controllers-k8s/runtime {{ .RuntimeVersion }}
	github.com/aws/aws-sdk-go {{ .AWSSDKGoVersion }}
	github.com/go-logr/logr v1.2.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.11.0
)
