module github.com/nginxinc/nginx-ingress-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/openshift/api v0.0.0-20201013121701-9d5ee23b507d
	github.com/prometheus/client_golang v1.9.0 // indirect
	golang.org/x/mod v0.4.0 // indirect
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.20.2
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	sigs.k8s.io/controller-runtime v0.8.3
)
