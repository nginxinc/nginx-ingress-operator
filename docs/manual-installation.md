# Manual installation

This will deploy the operator in the `nginx-ingress-operator-system` namespace.


1. Deploy the Operator and associated resources:
   1. Clone the `nginx-ingress-operator` repo and checkout the latest stable tag:
    ```
    git clone https://github.com/nginxinc/nginx-ingress-operator/
    cd nginx-ingress-operator/
    git checkout v0.5.1
    ```

   2. `Openshift` To deploy the Operator and associated resources to an OpenShift environment, run:
    ```
    make deploy IMG=registry.connect.redhat.com/nginx/nginx-ingress-operator:0.5.1
    ```

   3. Alternatively, to deploy the Operator and associated resources to all other environments:
    ```
    make deploy IMG=nginx/nginx-ingress-operator:0.5.1
    ```

2. Check that the Operator is running:
    ```
    kubectl get deployments -n nginx-ingress-operator-system

    NAME                                        READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-ingress-operator-controller-manager   1/1     1            1           15s
    ```
