# Changelog

### 0.0.7

FEATURES:

* [29](https://github.com/nginxinc/nginx-ingress-operator/pull/29) Add IngressClass support.
* [26](https://github.com/nginxinc/nginx-ingress-operator/pull/26) Add mTLS policy support.
* [25](https://github.com/nginxinc/nginx-ingress-operator/pull/25) Add JWT policy support.
* [21](https://github.com/nginxinc/nginx-ingress-operator/pull/21) Add latency metrics support.
* [18](https://github.com/nginxinc/nginx-ingress-operator/pull/18) Add support for policies in VS routes and VSR subroutes. Add RateLimit policy support

FIXES:

* [31](https://github.com/nginxinc/nginx-ingress-operator/pull/31) Add Status update for VS/VSR to RBAC.

KNOWS ISSUES:
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
