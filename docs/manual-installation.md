# Manual installation

### 1. Deploy the operator

This will deploy the operator in the `nginx-ingress-operator-system` namespace.


1. Deploy the Operator and associated resources:
   1. <Openshift> To deploy the Operator and associated resources to an OpenShift environment, run:
    ```
    make openshift-deploy IMG=registry.connect.redhat.com/nginx/nginx-ingress-operator:0.3.0
    ```

   2. To deploy the Operator and associated resources to all other environments:
    ```
    make deploy IMG=nginx/nginx-ingress-operator:0.3.0
    ```

2. Check that the Operator is running:
    ```
    kubectl get deployments -n nginx-ingress-operator-system   

    NAME                                        READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-ingress-operator-controller-manager   1/1     1            1           15s
    ```
