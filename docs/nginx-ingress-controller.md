# NginxIngressController Custom Resource

The `NginxIngressController` Custom Resource is the definition of a deployment of the Ingress Controller. 
With this Custom Resource, the NGINX Ingress Operator will be able to deploy and configure instances of the Ingress Controller in your cluster.  

## Configuration

There are several fields to configure the deployment of an Ingress Controller.

The following example shows the minimum configuration using only required fields:

```yaml
apiVersion: k8s.nginx.org/v1alpha1
kind: NginxIngressController
metadata:
  name: my-nginx-ingress-controller
  namespace: my-nginx-ingress
spec:
  type: deployment
  image:
    repository: nginx/nginx-ingress
    tag: edge
    pullPolicy: Always
  serviceType: NodePort
```

 The following example shows the usage of all fields (required and optional):
 
```yaml
 apiVersion: k8s.nginx.org/v1alpha1
 kind: NginxIngressController
 metadata:
   name: my-nginx-ingress-controller
   namespace: my-nginx-ingress
 spec:
   type: deployment
   nginxPlus: false
   image:
     repository: nginx/nginx-ingress
     tag: edge
     pullPolicy: Always
   replicas: 3
   serviceType: NodePort
   enableCRDs: true
   enableSnippets: false
   enablePreviewPolicies: false
   defaultSecret: my-nginx-ingress/default-secret
   ingressClass: my-nginx-ingress
   useIngressClassOnly: true
   watchNamespace: default
   healthStatus:
     enable: true
     uri: "/my-health"
   nginxDebug: true
   logLevel: 3
   nginxStatus:
     enable: true
     port: 9090
     allowCidrs: "127.0.0.1"
   enableLeaderElection: true
   wildcardTLS: my-nginx-ingress/wildcard-secret
   reportIngressStatus:
     enable: true
     externalService: my-nginx-ingress
   prometheus:
     enable: true
     port: 9114
   enableLatencyMetrics: false
   configMapData:
     error-log-level: debug
   enableTLSPassthrough: true
   globalConfiguration: my-nginx-ingress/nginx-configuration
   nginxReloadTimeout: 5000
   appProtect:
     enable: false
 ``` 
 
| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `type` | `string` | The type of the Ingress Controller installation - `deployment` or `daemonset`. | Yes |
| `nginxPlus` | `boolean` | Deploys the Ingress Controller for NGINX Plus. The default is `false` meaning the Ingress Controller will be deployed for NGINX OSS. | No |
| `image` | [image](#nginxingresscontrollerimage) | The image of the Ingress Controller. | Yes |
| `replicas` | `int` | The number of replicas of the Ingress Controller pod. The default is 1. Only applies if the `type` is set to deployment. | No |
| `defaultSecret` | `string` | The TLS Secret for TLS termination of the default server. The format is namespace/name. The secret must be of the type kubernetes.io/tls. If not specified, the operator will generate and deploy a TLS Secret with a self-signed certificate and key. | No |
| `serviceType` | `string` | The type of the Service for the Ingress Controller. Valid Service types are `NodePort` or `LoadBalancer`. | Yes |
| `enableCRDs` | `boolean` | Enables the use of NGINX Ingress Resource Definitions (VirtualServer and VirtualServerRoute). Default is `true`. | No |
| `enableSnippets` | `boolean` | Enable custom NGINX configuration snippets in VirtualServer, VirtualServerRoute and TransportServer resources. Requires `enableCRDs` set to `true`. | No |
| `enablePreviewPolicies` | `boolean` | Enables preview policies. Requires `enableCRDs` set to `true`. | No |
| `ingressClass` | `string` | A class of the Ingress controller. For Kubernetes >= 1.18, the Ingress controller only processes resources that belong to its class - i.e. have the "ingressClassName" field resource equal to the class. Additionally the Ingress Controller processes all the VirtualServer/VirtualServerRoute resources that do not have the "ingressClassName" field. For Kubernetes < 1.18, the Ingress Controller only processes resources that belong to its class - i.e have the annotation "kubernetes.io/ingress.class" (for Ingress resources) or field "ingressClassName" (for VirtualServer/VirtualServerRoute resources) equal to the class. Additionally, the Ingress Controller processes resources that do not have the class set, which can be disabled by setting `useIngressClassOnly` to `true`. Default is `nginx`. | No |
| `service` | [service](#nginxingresscontrollerservice) | The service of the Ingress Controller. | No |
| `useIngressClassOnly` | `boolean` | Ignore Ingress resources without the `"kubernetes.io/ingress.class"` annotation. For kubernetes versions >= 1.18 this flag will be IGNORED. | No |
| `watchNamespace` | `boolean` | Namespace to watch for Ingress resources. By default the Ingress controller watches all namespaces. | No |
| `healthStatus` | [healthStatus](#nginxingresscontrollerhealthstatus) | Adds a new location to the default server. The location responds with the 200 status code for any request. Useful for external health-checking of the Ingress Controller. | No |
| `nginxDebug` | `boolean` | Enable debugging for NGINX. Uses the nginx-debug binary. Requires `error-log-level: debug` in the configMapData. | No |
| `logLevel` | `int` | Log level for V logs. Format is `0 - 3` | No |
| `nginxStatus` | [nginxStatus](#nginxingresscontrollernginxstatus) | Configures NGINX stub_status, or the NGINX Plus API. | No |
| `reportIngressStatus` | [reportIngressStatus](#nginxingresscontrollerreportingressstatus) | Update the address field in the status of Ingresses resources. | No |
| `enableLeaderElection` | `boolean` | Enables Leader election to avoid multiple replicas of the controller reporting the status of Ingress resources – only one replica will report status. Default is `true`. | No |
| `wildcardTLS` | `string` | A Secret with a TLS certificate and key for TLS termination of every Ingress host for which TLS termination is enabled but the Secret is not specified. The secret must be of the type kubernetes.io/tls. If the argument is not set, for such Ingress hosts NGINX will break any attempt to establish a TLS connection. If the argument is set, but the Ingress controller is not able to fetch the Secret from Kubernetes API, the Ingress Controller will fail to start. Format is `namespace/name`. | No |
| `prometheus` | [prometheus](#nginxingresscontrollerprometheus) | Configures NGINX or NGINX Plus metrics in the Prometheus format. | No |
| `configMapData` | `map[string]string` | Initial values of the Ingress Controller ConfigMap. Check the [ConfigMap docs](https://docs.nginx.com/nginx-ingress-controller/configuration/global-configuration/configmap-resource/) for more information about possible values. | No |
| `globalConfiguration` | `string` | The GlobalConfiguration resource for global configuration of the Ingress Controller. Format is namespace/name. Requires `enableCRDs` set to `true`. | No |
| `enableTLSPassthrough` | `boolean` | Enable TLS Passthrough on port 443. Requires `enableCRDs` set to `true`. | No |
| `appprotect` | [appprotect](#nginxingresscontrollerappprotect) | App Protect support configuration. Requires `nginxPlus` set to `true`. | No |
| `nginxReloadTimeout` | `int`| Timeout in milliseconds which the Ingress Controller will wait for a successful NGINX reload after a change or at the initial start. (default is 4000. Default is 20000 instead if `enable-app-protect` is true) | No |

## NginxIngressController.Image

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `repository` | `string` | The repository of the image. | Yes |
| `tag` | `string` | The version of the image. | Yes |
| `pullPolicy` | `string` | The ImagePullPolicy of the image. Valid values are `Never`, `Always` or `IfNotPresent` | Yes |

## NginxIngressController.HealthStatus

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `enable` | `boolean` | Enable the HealthStatus. | Yes |
| `uri` | `string` | URI of the location. Default is `/nginx-health`. | No |

## NginxIngressController.NginxStatus

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `enable` | `boolean` | Enable the NginxStatus. | Yes |
| `port` | `int` | Set the port where the NGINX stub_status or the NGINX Plus API is exposed. Default is `8080`. Format is `1023 - 65535` | No |
| `allowCidrs` | `string` | Whitelist IPv4 IP/CIDR blocks to allow access to NGINX stub_status or the NGINX Plus API. Separate multiple IP/CIDR by commas. (default `127.0.0.1`) | No |

## NginxIngressController.Service

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `extraLabels` | `map[string]string` | Specifies extra labels of the service. | No |

## NginxIngressController.ReportIngressStatus

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `enable` | `boolean` | Enable reporting of the Ingress status. | Yes |
| `externalService` | `string` | Specifies the name of the service with the type LoadBalancer through which the Ingress controller pods are exposed externally. The external address of the service is used when reporting the status of Ingress resources. Note: if `serviceType` is `LoadBalancer`, the value of this field will be ignored, and the operator will use the name of the created LoadBalancer service instead. | No |
| `ingressLink` | `string` | Specifies the name of the IngressLink resource, which exposes the Ingress Controller pods via a BIG-IP system. The IP of the BIG-IP system is used when reporting the status of Ingress, VirtualServer and VirtualServerRoute resources. Requires `reportIngressStatus.enable` set to `true`. Note: If `serviceType` is `LoadBalancer` or `reportIngressStatus.externalService` is set, the value of this field will be ignored. | No |

## NginxIngressController.Prometheus

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `enable` | `boolean` | Enable Prometheus metrics. | Yes |
| `port` | `int` | Sets the port where the Prometheus metrics are exposed. Default is 9113. Format is `1023 - 65535`. | No |
| `enableLatencyMetrics` | `boolean` | Bucketed response times from when NGINX establishes a connection to an upstream server to when the last byte of the response body is received by NGINX. **Note** The metric for the upstream isn't available until traffic is sent to the upstream. Requires prometheus set to true | No |

## NginxIngressController.AppProtect

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `enable` | `boolean` | Enable App Protect. | Yes |
