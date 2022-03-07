# NGINX Ingress Controller Helm Chart

## Introduction

This chart deploys the NGINX Ingress controller in your Kubernetes cluster.

## Prerequisites

  - A [Kubernetes Version Supported by the Ingress Controller](https://docs.nginx.com/nginx-ingress-controller/technical-specifications/#supported-kubernetes-versions)
  - Helm 3.0+.
  - Git.
  - If you’d like to use NGINX Plus:
    - To pull from the F5 Container registry, configure a docker registry secret using your JWT token from the MyF5 portal by following the instructions from [here](https://docs.nginx.com/nginx-ingress-controller/installation/using-the-jwt-token-docker-secret). Make sure to specify the secret using `controller.serviceAccount.imagePullSecretName` parameter.
    - Alternatively, pull an Ingress controller image with NGINX Plus and push it to your private registry by following the instructions from [here](https://docs.nginx.com/nginx-ingress-controller/installation/pulling-ingress-controller-image).
    - Alternatively, you can build an Ingress controller image with NGINX Plus and push it to your private registry by following the instructions from [here](https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image).
    - Update the `controller.image.repository` field of the `values-plus.yaml` accordingly.
  - If you’d like to use App Protect Dos, please install App Protect Dos Arbitrator helm chart. Make sure to install in the same namespace as the NGINX Ingress Controller. Note that if you install multiple NGINX Ingress Controllers in the same namespace, they will need to share the same Arbitrator because it is not possible to install more than one Arbitrator in a single namespace.


## Getting the Chart Sources

This step is required if you're installing the chart using its sources. Additionally, the step is also required for managing the custom resource definitions (CRDs), which the Ingress Controller requires by default, or for upgrading/deleting the CRDs.

1. Clone the Ingress controller repo:
    ```console
    $ git clone https://github.com/nginxinc/kubernetes-ingress/
    ```
2. Change your working directory to /deployments/helm-chart:
    ```console
    $ cd kubernetes-ingress/deployments/helm-chart
    $ git checkout v2.1.0
    ```

## Adding the Helm Repository

This step is required if you're installing the chart via the helm repository.

```console
$ helm repo add nginx-stable https://helm.nginx.com/stable
$ helm repo update
```

## Installing the Chart

### Installing the CRDs

By default, the Ingress Controller requires a number of custom resource definitions (CRDs) installed in the cluster. The Helm client will install those CRDs. If the CRDs are not installed, the Ingress Controller pods will not become `Ready`.

If you do not use the custom resources that require those CRDs (which corresponds to `controller.enableCustomResources` set to `false` and `controller.appprotect.enable` set to `false` and `controller.appprotectdos.enable` set to `false`), the installation of the CRDs can be skipped by specifying `--skip-crds` for the helm install command.

### Installing via Helm Repository

To install the chart with the release name my-release (my-release is the name that you choose):

For NGINX:
```console
$ helm install my-release nginx-stable/nginx-ingress
```

For NGINX Plus: (assuming you have pushed the Ingress controller image `nginx-plus-ingress` to your private registry `myregistry.example.com`)
```console
$ helm install my-release nginx-stable/nginx-ingress --set controller.image.repository=myregistry.example.com/nginx-plus-ingress --set controller.nginxplus=true
```

**Note**: If you wish to use the experimental repository, replace `stable` with `edge` and add the `--devel` flag.

### Installing Using Chart Sources

To install the chart with the release name my-release (my-release is the name that you choose):

For NGINX:
```console
$ helm install my-release .
```

For NGINX Plus:
```console
$ helm install my-release -f values-plus.yaml .
```

**Note**: If you wish to use the experimental repository, replace the value in the `tag` field inside the yaml files with `edge`.

The command deploys the Ingress controller in your Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

When deploying the Ingress controller, make sure to use your own TLS certificate and key for the default server rather than the default pre-generated ones. Read the [Configuration](#Configuration) section below to see how to configure a TLS certificate and key for the default server. Note that the default server returns the Not Found page with the 404 status code for all requests for domains for which there are no Ingress rules defined.

## Upgrading the Chart

### Upgrading the CRDs

Helm does not upgrade the CRDs during a release upgrade. Before you upgrade a release, run the following command to upgrade the CRDs:

```console
$ kubectl apply -f crds/
```
> **Note**: The following warning is expected and can be ignored: `Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply`.

> **Note**: Make sure to check the [release notes](https://www.github.com/nginxinc/kubernetes-ingress/releases) for a new release for any special upgrade procedures.

### Upgrading the Release

To upgrade the release `my-release`:

#### Upgrade Using Chart Sources:

```console
$ helm upgrade my-release .
```

#### Upgrade via Helm Repository:

```console
$ helm upgrade my-release nginx-stable/nginx-ingress
```

## Uninstalling the Chart

### Uninstalling the Release

To uninstall/delete the release `my-release`:

```console
$ helm uninstall my-release
```
The command removes all the Kubernetes components associated with the release and deletes the release.

### Uninstalling the CRDs

Uninstalling the release does not remove the CRDs. To remove the CRDs, run:

```console
$ kubectl delete -f crds/
```
> **Note**: This command will delete all the corresponding custom resources in your cluster across all namespaces. Please ensure there are no custom resources that you want to keep and there are no other Ingress Controller releases running in the cluster.

## Running Multiple Ingress Controllers

If you are running multiple Ingress Controller releases in your cluster with enabled custom resources, the releases will share a single version of the CRDs. As a result, make sure that the Ingress Controller versions match the version of the CRDs. Additionally, when uninstalling a release, ensure that you don’t remove the CRDs until there are no other Ingress Controller releases running in the cluster.

See [running multiple ingress controllers](https://docs.nginx.com/nginx-ingress-controller/installation/running-multiple-ingress-controllers/) for more details.

## Configuration

The following tables lists the configurable parameters of the NGINX Ingress controller chart and their default values.

Parameter | Description | Default
--- | --- | ---
`controller.name` | The name of the Ingress controller daemonset or deployment. | Autogenerated
`controller.kind` | The kind of the Ingress controller installation - deployment or daemonset. | deployment
`controller.nginxplus` | Deploys the Ingress controller for NGINX Plus. | false
`controller.nginxReloadTimeout` | The timeout in milliseconds which the Ingress Controller will wait for a successful NGINX reload after a change or at the initial start. | 60000
`controller.hostNetwork` | Enables the Ingress controller pods to use the host's network namespace. | false
`controller.nginxDebug` | Enables debugging for NGINX. Uses the `nginx-debug` binary. Requires `error-log-level: debug` in the ConfigMap via `controller.config.entries`. | false
`controller.logLevel` | The log level of the Ingress Controller. | 1
`controller.image.repository` | The image repository of the Ingress controller. | nginx/nginx-ingress
`controller.image.tag` | The tag of the Ingress controller image. | 2.1.0
`controller.image.pullPolicy` | The pull policy for the Ingress controller image. | IfNotPresent
`controller.config.name` | The name of the ConfigMap used by the Ingress controller. | Autogenerated
`controller.config.annotations` | The annotations of the Ingress controller configmap. | {}
`controller.config.entries` | The entries of the ConfigMap for customizing NGINX configuration. See [ConfigMap resource docs](https://docs.nginx.com/nginx-ingress-controller/configuration/global-configuration/configmap-resource/) for the list of supported ConfigMap keys. | {}
`controller.customPorts` | A list of custom ports to expose on the NGINX ingress controller pod. Follows the conventional Kubernetes yaml syntax for container ports. | []
`controller.defaultTLS.cert` | The base64-encoded TLS certificate for the default HTTPS server. If not specified, a pre-generated self-signed certificate is used. **Note:** It is recommended that you specify your own certificate. | A pre-generated self-signed certificate.
`controller.defaultTLS.key` | The base64-encoded TLS key for the default HTTPS server. **Note:** If not specified, a pre-generated key is used. It is recommended that you specify your own key. | A pre-generated key.
`controller.defaultTLS.secret` | The secret with a TLS certificate and key for the default HTTPS server. The value must follow the following format: `<namespace>/<name>`. Used as an alternative to specifying a certificate and key using `controller.defaultTLS.cert` and `controller.defaultTLS.key` parameters. | None
`controller.wildcardTLS.cert` | The base64-encoded TLS certificate for every Ingress/VirtualServer host that has TLS enabled but no secret specified. If the parameter is not set, for such Ingress/VirtualServer hosts NGINX will break any attempt to establish a TLS connection. | None
`controller.wildcardTLS.key` | The base64-encoded TLS key for every Ingress/VirtualServer host that has TLS enabled but no secret specified. If the parameter is not set, for such Ingress/VirtualServer hosts NGINX will break any attempt to establish a TLS connection. | None
`controller.wildcardTLS.secret` | The secret with a TLS certificate and key for every Ingress/VirtualServer host that has TLS enabled but no secret specified. The value must follow the following format: `<namespace>/<name>`. Used as an alternative to specifying a certificate and key using `controller.wildcardTLS.cert` and `controller.wildcardTLS.key` parameters. | None
`controller.nodeSelector` | The node selector for pod assignment for the Ingress controller pods. | {}
`controller.terminationGracePeriodSeconds` | The termination grace period of the Ingress controller pod. | 30
`controller.tolerations` | The tolerations of the Ingress controller pods. | []
`controller.affinity` | The affinity of the Ingress controller pods. | {}
`controller.volumes` | The volumes of the Ingress controller pods. | []
`controller.volumeMounts` | The volumeMounts of the Ingress controller pods. | []
`controller.initContainers` | InitContainers for the Ingress controller pods. | []
`controller.resources` | The resources of the Ingress controller pods. | {}
`controller.replicaCount` | The number of replicas of the Ingress controller deployment. | 1
`controller.ingressClass` | A class of the Ingress controller. An IngressClass resource with the name equal to the class must be deployed. Otherwise, the Ingress Controller will fail to start. The Ingress controller only processes resources that belong to its class - i.e. have the "ingressClassName" field resource equal to the class. The Ingress Controller processes all the VirtualServer/VirtualServerRoute/TransportServer resources that do not have the "ingressClassName" field for all versions of kubernetes. | nginx
`controller.setAsDefaultIngress` | New Ingresses without an `"ingressClassName"` field specified will be assigned the class specified in `controller.ingressClass`. | false
`controller.watchNamespace` | Namespace to watch for Ingress resources. By default the Ingress controller watches all namespaces. | ""
`controller.enableCustomResources` | Enable the custom resources. | true
`controller.enablePreviewPolicies` | Enable preview policies. | false
`controller.enableTLSPassthrough` | Enable TLS Passthrough on port 443. Requires `controller.enableCustomResources`. | false
`controller.globalConfiguration.create` | Creates the GlobalConfiguration custom resource. Requires `controller.enableCustomResources`. | false
`controller.globalConfiguration.spec` | The spec of the GlobalConfiguration for defining the global configuration parameters of the Ingress Controller. | {}
`controller.enableSnippets` | Enable custom NGINX configuration snippets in Ingress, VirtualServer, VirtualServerRoute and TransportServer resources. | false
`controller.healthStatus` | Add a location "/nginx-health" to the default server. The location responds with the 200 status code for any request. Useful for external health-checking of the Ingress controller. | false
`controller.healthStatusURI` | Sets the URI of health status location in the default server. Requires `controller.healthStatus`. | "/nginx-health"
`controller.nginxStatus.enable` | Enable the NGINX stub_status, or the NGINX Plus API. | true
`controller.nginxStatus.port` | Set the port where the NGINX stub_status or the NGINX Plus API is exposed. | 8080
`controller.nginxStatus.allowCidrs` | Add IPv4 IP/CIDR blocks to the allow list for NGINX stub_status or the NGINX Plus API. Separate multiple IP/CIDR by commas. | 127.0.0.1
`controller.priorityClassName` | The PriorityClass of the Ingress controller pods. | None
`controller.service.create` | Creates a service to expose the Ingress controller pods. | true
`controller.service.type` | The type of service to create for the Ingress controller. | LoadBalancer
`controller.service.externalTrafficPolicy` | The externalTrafficPolicy of the service. The value Local preserves the client source IP. | Local
`controller.service.annotations` | The annotations of the Ingress controller service. | {}
`controller.service.extraLabels` | The extra labels of the service. | {}
`controller.service.loadBalancerIP` | The static IP address for the load balancer. Requires `controller.service.type` set to `LoadBalancer`. The cloud provider must support this feature. | ""
`controller.service.externalIPs` | The list of external IPs for the Ingress controller service. | []
`controller.service.loadBalancerSourceRanges` | The IP ranges (CIDR) that are allowed to access the load balancer. Requires `controller.service.type` set to `LoadBalancer`. The cloud provider must support this feature. | []
`controller.service.name` | The name of the service. | Autogenerated
`controller.service.customPorts` | A list of custom ports to expose through the Ingress controller service. Follows the conventional Kubernetes yaml syntax for service ports. | []
`controller.service.httpPort.enable` | Enables the HTTP port for the Ingress controller service. | true
`controller.service.httpPort.port` | The HTTP port of the Ingress controller service. | 80
`controller.service.httpPort.nodePort` | The custom NodePort for the HTTP port. Requires `controller.service.type` set to `NodePort`. | ""
`controller.service.httpPort.targetPort` | The target port of the HTTP port of the Ingress controller service. | 80
`controller.service.httpsPort.enable` | Enables the HTTPS port for the Ingress controller service. | true
`controller.service.httpsPort.port` | The HTTPS port of the Ingress controller service. | 443
`controller.service.httpsPort.nodePort` | The custom NodePort for the HTTPS port. Requires `controller.service.type` set to `NodePort`.  | ""
`controller.service.httpsPort.targetPort` | The target port of the HTTPS port of the Ingress controller service. | 443
`controller.serviceAccount.name` | The name of the service account of the Ingress controller pods. Used for RBAC. | Autogenerated
`controller.serviceAccount.imagePullSecretName` | The name of the secret containing docker registry credentials. Secret must exist in the same namespace as the helm release. | ""
`controller.reportIngressStatus.enable` | Updates the address field in the status of Ingress resources with an external address of the Ingress controller. You must also specify the source of the external address either through an external service via `controller.reportIngressStatus.externalService`, `controller.reportIngressStatus.ingressLink` or the `external-status-address` entry in the ConfigMap via `controller.config.entries`. **Note:** `controller.config.entries.external-status-address` takes precedence over the others. | true
`controller.reportIngressStatus.externalService` | Specifies the name of the service with the type LoadBalancer through which the Ingress controller is exposed externally. The external address of the service is used when reporting the status of Ingress, VirtualServer and VirtualServerRoute resources. `controller.reportIngressStatus.enable` must be set to `true`. The default is autogenerated and enabled when `controller.service.create` is set to `true` and `controller.service.type` is set to `LoadBalancer`. | Autogenerated
`controller.reportIngressStatus.ingressLink` |  Specifies the name of the IngressLink resource, which exposes the Ingress Controller pods via a BIG-IP system. The IP of the BIG-IP system is used when reporting the status of Ingress, VirtualServer and VirtualServerRoute resources. `controller.reportIngressStatus.enable` must be set to `true`. | ""
`controller.reportIngressStatus.enableLeaderElection` | Enable Leader election to avoid multiple replicas of the controller reporting the status of Ingress resources. `controller.reportIngressStatus.enable` must be set to `true`. | true
`controller.reportIngressStatus.leaderElectionLockName` | Specifies the name of the ConfigMap, within the same namespace as the controller, used as the lock for leader election. controller.reportIngressStatus.enableLeaderElection must be set to true. | Autogenerated
`controller.reportIngressStatus.annotations` | The annotations of the leader election configmap. | {}
`controller.pod.annotations` | The annotations of the Ingress Controller pod. | {}
`controller.pod.extraLabels` | The additional extra labels of the Ingress Controller pod. | {}
`controller.appprotect.enable` | Enables the App Protect module in the Ingress Controller. | false
`controller.appprotectdos.enable` | Enables the App Protect Dos module in the Ingress Controller. | false
`controller.appprotectdos.debug` | Enable debugging for App Protect Dos. | false
`controller.appprotectdos.maxDaemons` | Max number of ADMD instances. | 1
`controller.appprotectdos.maxWorkers` | Max number of nginx processes to support. | Number of CPU cores in the machine
`controller.appprotectdos.memory` | RAM memory size to consume in MB. | 50% of free RAM in the container or 80MB, the smaller
`controller.readyStatus.enable` | Enables the readiness endpoint `"/nginx-ready"`. The endpoint returns a success code when NGINX has loaded all the config after the startup. This also configures a readiness probe for the Ingress Controller pods that uses the readiness endpoint. | true
`controller.readyStatus.port` | The HTTP port for the readiness endpoint. | 8081
`controller.enableLatencyMetrics` |  Enable collection of latency metrics for upstreams. Requires `prometheus.create`. | false
`rbac.create` | Configures RBAC. | true
`prometheus.create` | Expose NGINX or NGINX Plus metrics in the Prometheus format. | false
`prometheus.port` | Configures the port to scrape the metrics. | 9113
`prometheus.scheme` | Configures the HTTP scheme to use for connections to the Prometheus endpoint. | http
`prometheus.secret` | The namespace / name of a Kubernetes TLS Secret. If specified, this secret is used to secure the Prometheus endpoint with TLS connections. | ""


## Notes
* The values-icp.yaml file is used for deploying the Ingress controller on IBM Cloud Private. See the [blog post](https://www.nginx.com/blog/nginx-ingress-controller-ibm-cloud-private/) for more details.
