# Changelog

### 0.2.0

An automatically generated list of changes can be found on Github at: [0.2.0 Release](https://github.com/nginxinc/nginx-ingress-operator/releases/tag/v0.2.0)

### 0.1.0

FEATURES:

* [56](https://github.com/nginxinc/nginx-ingress-operator/pull/56) Graduate Policies. Add enablePreviewPolicy flag support.
* [55](https://github.com/nginxinc/nginx-ingress-operator/pull/55) Add AppProtect User Defined Signatures support.
* [39](https://github.com/nginxinc/nginx-ingress-operator/pull/39) Update secret type of default secret to TLS.

FIXES:

* [71](https://github.com/nginxinc/nginx-ingress-operator/pull/71) Fix replicas and service to be optional fields.
* [70](https://github.com/nginxinc/nginx-ingress-operator/pull/70) Make enableCRDs optional.
* [66](https://github.com/nginxinc/nginx-ingress-operator/pull/66) Fix Service to be an optional field. Add support for updating ExtraLabels.
* [65](https://github.com/nginxinc/nginx-ingress-operator/pull/65) Fix SCC resource to only affect KIC pods.

DOCUMENTATION:

* [54](https://github.com/nginxinc/nginx-ingress-operator/pull/54) Update IC compatibility in changelog.

KNOWN ISSUES:

* The Operator doesn't automatically remove IngressClasses created by [29](https://github.com/nginxinc/nginx-ingress-operator/pull/29).

COMPATIBILITY:

- NGINX Ingress Controller 1.10.x
- Openshift 4.5 or newer.

UPGRADE INSTRUCTIONS:

1. Remove the existing Policy CRD: `kubectl delete crd policies.k8s.nginx.org`
**Please note that deletion of the `policies.k8s.nginx.org` CRD will result in all instances of that CRD being deleted too. Ensure to back up any important Custom Resource instances first!**
1. Delete the existing SCC: `kubectl delete scc nginx-ingress-scc`
1. Upgrade the operator to version 0.1.0.
1. If the defaultSecret field is not set in your `nginxingresscontrollers.k8s.nginx.org` resource (or resources):
    1. Remove the generated default secret. For example: `kubectl delete secret -n my-nginx-ingress my-nginx-ingress-controller`
    1. Wait until the operator regenerates the secret. The old secret was of the type `Opaque`. The new secret will be of the type `kubernetes.io/tls`.
1. Alternatively, if the defaultSecret is set to some secret, make sure it is of the type `kubernetes.io/tls`. If not, recreate the secret with the type `kubernetes.io/tls`.
1. If the wildcardTLS is set to some secret, make sure it is of the type `kubernetes.io/tls`. If not, recreate the secret with the type `kubernetes.io/tls`.
1. Ensure that the TLS secrets referenced by Ingress, VirtualServer and Policy resources are of the type `kubernetes.io/tls`, JWT secrets are of the type `nginx.org/jwt` and CA secrets are of the type `nginx.org/ca`. To avoid potential disruption of client traffic, instead of recreating the secrets, create new secrets with the correct type and update the Ingress, VirtualServer and Policy resources to use the new secrets.
1. Update any existing instances of the `nginxingresscontrollers.k8s.nginx.org` Custom Resource to use an NGINX Ingress Controller 1.10.x image.

**Note**: Steps 4-8 are required because Version 1.10.0 of the Ingress Controller added a requirement for secrets to be one of the following types: `kubernetes.io/tls` for TLS secrets; `nginx.org/jwk` for JWK secrets; or `nginx.org/ca` for CA secrets. Please see the section UPDATING SECRETS in https://docs.nginx.com/nginx-ingress-controller/releases/#nginx-ingress-controller-1-10-0 for more details.

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
