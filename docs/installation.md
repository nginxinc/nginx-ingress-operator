# Installation

## Prerequisites
* Go >= 1.13
* Docker

### 1. Build the image

1. Build the operator image
    ```
    make build
    ```

1. Push the image to your private repository.
    ```
    docker push nginx-ingress-operator
    ```

### 2. Deploy the operator

This will deploy the operator in the `default` namespace.

1. Deploy the NginxIngressController Custom Resource Definition:
    ```
    kubectl create -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml
    ```

1. Deploy the ServiceAccount:
    ```
    kubectl create -f deploy/service_account.yaml
    ```
   
1. Deploy the Role:
    ```
    kubectl create -f deploy/role.yaml
    ```

1. Deploy the RoleBinding:
    ```
    kubectl create -f deploy/role_binding.yaml
    ```

1. Deploy the Operator:

    **Note**: Update the `image` field with your previously built image if necessary.

    ```
    kubectl create -f deploy/operator.yaml
    ```

1. Check that the operator is running:
    ```
    kubectl get deployment nginx-ingress-operator
        
    NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-ingress-operator   1/1     1            1           15s
    ```    
   
Check the [examples](../examples) to deploy the NGINX Ingress Controller using the operator.
