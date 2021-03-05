# Changelog

### 0.1.0

Features

* Add IngressLink support (#58) @Dean-Coakley
* Add AppProtect User Defined Signatures support (#55) @Dean-Coakley
* Graduate Policies. Add enablePreviewPolicy flag support. (#56) @Dean-Coakley
* Update secret type of default secret to TLS (#39) @pleshakov

Bug Fixes

* Fix replicas and service to be optional fields (#71) @Dean-Coakley
* Make enableCRDs optional (#70) @Dean-Coakley
* Fix Service to be an optional field. Add support for updating ExtraLabels. (#66) @Dean-Coakley
* Fix SCC resource to only affect KIC pods (#65) @Dean-Coakley

Documentation

* Update IC compatibility in changelog (#54) @pleshakov

Maintenance

* Bump golangci/golangci-lint-action from v2.4.0 to v2.5.1 (#67) @dependabot
* Add release-drafter (#64) @lucacome
* Change dependabot interval to weekly (#63) @lucacome
* Bump actions/cache from v2 to v2.1.4 (#53) @dependabot
* Bump golangci/golangci-lint-action from v2 to v2.4.0 (#59) @dependabot
* Bump github.com/google/go-cmp from 0.4.0 to 0.5.4 (#49) @dependabot
* Add dependabot (#45) @lucacome
* Update CRDs, CSVs and Makefile (#36) @lucacome

Compatibility

- NGINX Ingress Controller 1.10.x
- Openshift 4.5 or newer.

Upgrade Instructions

1. Remove existing policy CRD: `kubectl delete crds policies.k8s.nginx.org`
  **Please note that deletion of the policies.k8s.nginx.org CRD will result in all instances of that CRD being deleted too. Ensure to back up any important Custom Resource instances first!**
2. Delete existing SCC: `kubectl delete scc nginx-ingress-scc`
3. Deploy new operator.
4. Update any existing instances of the nginxingresscontrollers.k8s.nginx.org Custom Resource to use a KIC 1.10 image.

### 0.0.7

FEATURES:

* [29](https://github.com/nginxinc/nginx-ingress-operator/pull/29) Add IngressClass support.
* [26](https://github.com/nginxinc/nginx-ingress-operator/pull/26) Add mTLS policy support.
* [25](https://github.com/nginxinc/nginx-ingress-operator/pull/25) Add JWT policy support.
* [21](https://github.com/nginxinc/nginx-ingress-operator/pull/21) Add latency metrics support.
* [18](https://github.com/nginxinc/nginx-ingress-operator/pull/18) Add support for policies in VS routes and VSR subroutes. Add RateLimit policy support

FIXES:

* [31](https://github.com/nginxinc/nginx-ingress-operator/pull/31) Add Status update for VS/VSR to RBAC.

KNOWN ISSUES:
* The Operator doesn't automatically remove IngressClasses created by [29](https://github.com/nginxinc/nginx-ingress-operator/pull/29)

COMPATIBILITY:

* NGINX Ingress Controller 1.9.x.
* Openshift 4.5 or newer.

### 0.0.6

FEATURES:

* [13](https://github.com/nginxinc/nginx-ingress-operator/pull/13) Add support for App Protect.
* [11](https://github.com/nginxinc/nginx-ingress-operator/pull/11) Add enableSnippets cli argument support.

IMPROVEMENTS:
* [15](https://github.com/nginxinc/nginx-ingress-operator/pull/15) Downgrade operator-sdk to 0.17.
* [14](https://github.com/nginxinc/nginx-ingress-operator/pull/14) Add KIC supported versions to README.
* [12](https://github.com/nginxinc/nginx-ingress-operator/pull/12) Make operator install KIC CRDs from manifests.
* [10](https://github.com/nginxinc/nginx-ingress-operator/pull/10) Update operator-sdk to 0.18.
* [8](https://github.com/nginxinc/nginx-ingress-operator/pull/8) Update go to 1.14.
* [7](https://github.com/nginxinc/nginx-ingress-operator/pull/7) Update makefile to include all manifests.

COMPATIBILITY:

* NGINX Ingress Controller 1.8.x.
* Openshift 4.3 or newer.

### 0.0.4

FEATURES:

* [4](https://github.com/nginxinc/nginx-ingress-operator/pull/4) Add new CRDs for NGINX Ingress Controller 1.7.0
* [5](https://github.com/nginxinc/nginx-ingress-operator/pull/5) Make NGINX Ingress Operator RedHat certified. Learn more about certified operators for Openshift [here](https://connect.redhat.com/en/partner-with-us/red-hat-openshift-operator-certification).

COMPATIBILITY:

* NGINX Ingress Controller 1.7.x.
* Openshift 4.3 or newer.
