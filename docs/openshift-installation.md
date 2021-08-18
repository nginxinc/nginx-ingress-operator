# Installation in an Openshift cluster the OLM

This installation method is the recommended way for Openshift users. **Note**: Openshift version must be 4.2 or higher.

**Note: The `nginx-ingress-operator` supports `Basic Install` only - we do not support auto-updates. When you are installing the Operator using the OLM, the auto-update feature should be disabled to avoid breaking changes being auto-applied. In OpenShift, this can be done by setting the `Approval Strategy` to `Manual`. Please see the [Operator SDK docs](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/) for more details on the Operator Capability Levels.**

The NGINX Ingress Operator is a [RedHat certified Operator](https://connect.redhat.com/en/partner-with-us/red-hat-openshift-operator-certification).

1. In the Openshift dashboard, click `Operators` > `Operator Hub` in the left menu and use the search box to type `nginx ingress`:
![alt text](./images/openshift1.png "Operators")
1. Click the `NGINX Ingress Operator` and click `Install`:
![alt text](./images/openshift2.png "NGINX Ingress Operator")
1. Click `Subscribe`:
![alt text](./images/openshift3.png "NGINX Ingress Operator Install")

Openshift will install the NGINX Ingress Operator:

![alt text](./images/openshift4.png "NGINX Ingress Operator Subscribe")

You can now deploy the NGINX Ingress Controller instances following the [examples](../examples).
