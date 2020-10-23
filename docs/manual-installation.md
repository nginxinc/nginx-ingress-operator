# Manual installation

### 1. Deploy the operator

This will deploy the operator in the `default` namespace.

1. Deploy the NginxIngressController Custom Resource Definition:
    ```
    kubectl apply -f deploy/crds/k8s.nginx.org_nginxingresscontrollers_crd.yaml
    ```

1. Deploy the ServiceAccount:
    ```
    kubectl apply -f deploy/service_account.yaml
    ```

1. Deploy the Role:
    ```
    kubectl apply -f deploy/role.yaml
    ```

1. Deploy the RoleBinding:
    ```
    kubectl apply -f deploy/role_binding.yaml
    ```

1. Deploy the Operator:
    ```
    kubectl apply -f deploy/operator.yaml
    ```

1. Check that the Operator is running:
    ```
    kubectl get deployment nginx-ingress-operator

    NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-ingress-operator   1/1     1            1           15s
    ```