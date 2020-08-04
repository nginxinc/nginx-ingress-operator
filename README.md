# NGINX Ingress Operator

The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more [NGINX/NGINX Plus Ingress Controllers](https://github.com/nginxinc/kubernetes-ingress) which in turn handle Ingress traffic for applications running in a cluster.

Learn more about operators in the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

To install a specific version of the NGINX Ingress Controller with the operator, a specific version of the NGINX Ingress Operator is required.

The following table shows the relation between the versions of the two projects:

| NGINX Ingress Controller | NGINX Ingress Operator |
| --- | --- |
| < 1.7.0 | N/A |
| 1.7.x | 0.0.4 |
| 1.8.x | 0.0.6 |

Note: The NGINX Ingress Operator works only for NGINX Ingress Controller versions after `1.7.0`.

## Getting Started

1. Install the NGINX Ingress Operator. See [docs](./docs/installation.md).
1. Deploy a new NGINX Ingress Controller using the [NginxIngressController](docs/nginx-ingress-controller.md) Custom Resource:
    * For an NGINX installation see the [NGINX example](./examples/deployment-oss-min).
    * For an NGINX Plus installation see the [NGINX Plus example](./examples/deployment-plus-min).

## NGINX Ingress Operator Releases
We publish NGINX Ingress Operator releases on GitHub. See our [releases page](https://github.com/nginxinc/nginx-ingress-operator/releases).

The latest stable release is [0.0.6](https://github.com/nginxinc/nginx-ingress-operator/releases/tag/v0.0.6). For production use, we recommend that you choose the latest stable release.

## Development

It is possible to run the operator in your local machine. This is useful for testing or during development.

### Run Operator locally

1. Have access to a Kubernetes/Openshift cluster.
1. Apply the latest CRD:
    ```
    kubectl apply -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml
    ```
1. Apply the NGINX Ingress Controller CRDs:
    ```
    kubectl apply -f build/kic_crds
    ```
1. Run `make run-local`.

The operator will run in your local machine but will be communicating with the cluster. The operator will only watch the `default` namespace when deployed locally.

### Update CRD

If any change is made in the CRD in the go code, run the following commands to update the changes in the CRD yaml:

1. `make generate-crds`
1. Apply the new CRD definition again in your cluster `kubectl apply -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml`.

### Run tests

Run `make test` to run unit tests locally.

## Contributing

If you'd like to contribute to the project, please read our [Contributing](./CONTRIBUTING.md) guide.
