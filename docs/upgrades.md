# Upgrade - 0.2.0 to 0.3.0

Release 0.3.0 includes a major upgrade of the Operator-SDK which has resulted in a number of changes in the layout of the project
(see [the operator-sdk docs](https://sdk.operatorframework.io/docs/building-operators/golang/migration/) for more information).

## OLM upgrade - 0.2.0 to 0.3.0

**Note: The `nginx-ingress-operator` supports basic install only - we do not support auto-updates. When you are installing the Operator using the OLM, the auto-update feature should be disabled to avoid breaking changes being auto-applied.**

### 1. Uninstall the existing 0.2.0 operator

Uninstall the operator using the web console - see [the OCP documentation for details](https://access.redhat.com/documentation/en-us/openshift_container_platform/4.6/pdf/operators/OpenShift_Container_Platform-4.6-Operators-en-US.pdf).

### 2. Reinstall the latest version of the operator

Install the latest version of the Operator following the steps outlined in [OpenShift installation doc](./openshift-installation.md).

### 3. Upgrade the existing ingress controller deployments

Upgrade to the latest 1.12.0 Ingress Controller image - see the release notes [here](https://docs.nginx.com/nginx-ingress-controller/releases/#nginx-ingress-controller-1-12-0)

## Manual upgrade - 0.2.0 to 0.3.0

### 1. Deploy the new operator

Deploy the operator following the steps outlined in [manual installation doc](./manual-installation.md).

### 2. Cleanup the existing operator

Uninstall the existing operator deployment:
   
1. Checkout the previous version of the nginx-ingress-operator [0.3.0](https://github.com/nginxinc/nginx-ingress-operator/releases/tag/v0.3.0).
1. Uninstall the resources by running the following commands (be sure to edit files to suit your environment, if required):
    ```
    kubectl delete -f deploy/operator.yaml
    kubectl delete -f deploy/role_binding.yaml
    kubectl delete -f deploy/role.yaml
    kubectl delete -f deploy/service_account.yaml
    ```

### 3. Upgrade the existing ingress controller deployments

Upgrade to the latest 1.12.0 Ingress Controller image - see the release notes [here](https://docs.nginx.com/nginx-ingress-controller/releases/#nginx-ingress-controller-1-12-0)
