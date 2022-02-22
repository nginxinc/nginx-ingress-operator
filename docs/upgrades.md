# Upgrade - 0.5.0 to 1.0.0

Release 1.0.0 includes a backward incompatible change from version 0.5.0 as we have moved from a Go based operator to a Helm based operator.

## OLM upgrade - 0.5.0 to 1.0.0

**Note: The `nginx-ingress-operator` supports `Basic Install` only - we do not support auto-updates. When you are installing the Operator using the OLM, the auto-update feature should be disabled to avoid breaking changes being auto-applied. In OpenShift, this can be done by setting the `Approval Strategy` to `Manual`. Please see the [Operator SDK docs](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/) for more details on the Operator Capability Levels.**
1. Upgrade CRDs
2. Uninstall Go operator -> this will also remove any instances of the NginxIngressController, but not any dependent objects (ingresses, VSs, etc)
3. Remove the nginx-ingress ingressClass `k delete ingressclass/nginx`
4. Install new operator 
5. Deploy common resources (scc, default server secret, ns, etc). Note: service account and ingress class should be deployed separately if deploying multiple ICs in same namespace. This is because only one of the ICs in a namespace will be assigned "ownership" of these resources.
6. Re-create ingress controllers (note: multi IC rules) using the new Operator. Be sure to use the same configuration as the previous deployments (ingress class name, namespaces etc). They will pick up all deployed resources.

### 0. Upgrade the existing NIC crds

Navigate [here]() and run ` kubectl apply -f crds/`

### 1. Uninstall the existing 0.5.0 operator, the nginx ingress controller CRD, and the ingressClass

Uninstall the operator using the web console - see [the OCP documentation for details](https://access.redhat.com/documentation/en-us/openshift_container_platform/4.9/pdf/operators/OpenShift_Container_Platform-4.9-Operators-en-US.pdf). 

Next uninstall the NIC CRD. This will remove any instances of the NginxIngressController, but not any dependent objects (ingresses, VSs, etc).

Finally, remove the nginx-ingress ingressClass `k delete ingressclass/nginx`.

### 2. Install the latest version of the operator

Install the latest version of the Operator following the steps outlined in [OpenShift installation doc](./openshift-installation.md).

### 3. Deploy new ingress controller deployments

Use the new Nginx Ingress Operator installation to deploy Nginx Ingress Controller - see the release notes [here](https://docs.nginx.com/nginx-ingress-controller/releases/#nginx-ingress-controller-2-2-0) and a guide to the Helm configuration parameters [here](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/#configuration)

## Manual upgrade - 0.5.0 to 1.0.0

### 1. Deploy the new operator

Deploy the operator following the steps outlined in [manual installation doc](./manual-installation.md).

### 2. Cleanup the existing operator

Uninstall the existing operator deployment:
   
1. Checkout the previous version of the nginx-ingress-operator [0.5.0](https://github.com/nginxinc/nginx-ingress-operator/releases/tag/v0.5.0).
2. Uninstall the resources by running the following command:
    ```
    make undeploy
    ```

### 3. Upgrade the existing ingress controller deployments

Upgrade to the latest 2.2.0 Ingress Controller image - see the release notes [here](https://docs.nginx.com/nginx-ingress-controller/releases/#nginx-ingress-controller-2-2-0)
