module github.com/nginxinc/nginx-ingress-operator

go 1.16

require (
	github.com/go-logr/logr v0.1.0
	github.com/golangci/golangci-lint v1.39.0
	github.com/google/go-cmp v0.5.5
	github.com/openshift/api v0.0.0-20200701144905-de5b010b2b38
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.18.9
	k8s.io/apiextensions-apiserver v0.18.9
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.3
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.9 // Required by prometheus-operator
)
