# Installation in an Openshift cluster the OLM

This installation method is the recommended way for Openshift users. **Note**: Openshift version must be 4.2 or higher.

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