# NGINX Ingress Operator

The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more [NGINX/NGINX Plus Ingress Controllers](https://github.com/nginxinc/kubernetes-ingress) which in turn handle Ingress traffic for applications running in a cluster.

Learn more about operators in the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

## Getting Started

1. Install the NGINX Ingress Operator. See [docs](./docs/installation.md).
1. Deploy a new NGINX Ingress Controller using the [NginxIngressController](docs/nginx-ingress-controller.md) Custom Resource:
    * For an NGINX installation see the [NGINX example](./examples/deployment-oss-min).
    * For an NGINX Plus installation see the [NGINX Plus example](./examples/deployment-plus-min).

## Development

It is possible to run the operator in your local machine. This is useful for testing or during development.

### Run Operator locally

1. Have access to a Kubernetes/Openshift cluster.
1. Apply the latest CRD:
    ```
    kubectl apply -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml
    ```
1. Run `make run-local`.

The operator will run in your local machine but will be communicating with the cluster.

### Update CRD

If any change is made in the CRD in the go code, run the following commands to update the changes in the CRD yaml:

1. `make generate-crds`
1. Apply the new CRD definition again in your cluster `kubectl apply -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml`.

### Run tests

Run `make test` to run unit tests locally.

## Contributing

If you'd like to contribute to the project, please read our [Contributing](./CONTRIBUTING.md) guide.