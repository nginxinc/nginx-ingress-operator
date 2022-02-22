[![Continuous Integration](https://github.com/nginxinc/nginx-ingress-operator/workflows/Continuous%20Integration/badge.svg)](https://github.com/nginxinc/nginx-ingress-operator/actions) [![FOSSA Status](https://app.fossa.com/api/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fnginx-ingress-operator.svg?type=shield)](https://app.fossa.com/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fnginx-ingress-operator?ref=badge_shield)

# NGINX Ingress Operator

The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more [NGINX/NGINX Plus Ingress Controllers](https://github.com/nginxinc/kubernetes-ingress) which in turn handle Ingress traffic for applications running in a cluster.

Learn more about operators in the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

To install a specific version of the NGINX Ingress Controller with the operator, a specific version of the NGINX Ingress Operator is required.

Up until version 0.5.0, this Operator was Go based. Version 1.0.0 marks an uncompatible upgrade as this release swtiched the Operator to being Helm-based, built from the [NGINX Ingress Controller Helm chart](http://helm.nginx.com/#nginx-ingress-controller). The configuration for the Helm chart can be seen in the [NGINX Ingress Controller documentation](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/#configuration).

The following table shows the relation between the versions of the two projects:

| NGINX Ingress Controller | NGINX Ingress Operator |
| --- | --- |
| 2.2.x | 1.0.0 |
| 2.1.x | 0.5.0 |
| 2.0.x | 0.4.0 |
| 1.12.x | 0.3.0 |
| 1.11.x | 0.2.0 |
| 1.10.x | 0.1.0 |
| 1.9.x | 0.0.7 |
| 1.8.x | 0.0.6 |
| 1.7.x | 0.0.4 |
| < 1.7.0 | N/A |

Note: The NGINX Ingress Operator works only for NGINX Ingress Controller versions after `1.7.0`.

## Getting Started

1. Install the NGINX Ingress Operator. See [docs](./docs/installation.md).
   <br> NOTE: To use TransportServers as part of your NGINX Ingress Controller configuration, a GlobalConfiguration resource must be created *before* starting the Operator - [see the notes](./examples/deployment-oss-min/README.md#TransportServers)
2. Create a default server secret on the cluster - an example yaml for this can be found in the [NGINX Ingress Controller GitHub repo](https://github.com/nginxinc/kubernetes-ingress/blob/master/deployments/common/default-server-secret.yaml)
3. Deploy a new NGINX Ingress Controller using the [NginxIngress](./config/samples/charts_v1alpha1_nginxingress.yaml) Custom Resource:
    * Use the name of the default server secret created above for `controller.defaultTLS.secret` field (needs to be in the form `namespace/name`)
    * If using NGINX Plus:
      * Set the `controller.nginxPlus` to true
      * Set the `controller.image.repository` and `controller.image.tag` to the appropriate values
      * Set the `controller.serviceAccount.imagePullSecretName` if applicable
    * For full configuration details see the Helm documentation [here](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/#configuration).

## Upgrades

See [upgrade docs](./docs/upgrades)

## NGINX Ingress Operator Releases
We publish NGINX Ingress Operator releases on GitHub. See our [releases page](https://github.com/nginxinc/nginx-ingress-operator/releases).

The latest stable release is [1.0.0](https://github.com/nginxinc/nginx-ingress-operator/releases/tag/v1.0.0). For production use, we recommend that you choose the latest stable release.

## Development

It is possible to run the operator in your local machine. This is useful for testing or during development.

### Run Operator locally

1. Have access to a Kubernetes/Openshift cluster.
1. Apply the IC CRD:
   ```
    make install
    ```
2. Run `make run`.

The operator will run in your local machine but will be communicating with the cluster.

## Contributing

If you'd like to contribute to the project, please read our [Contributing](./CONTRIBUTING.md) guide.
