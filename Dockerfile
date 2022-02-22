FROM registry.redhat.io/openshift4/ose-helm-operator:v4.9

ARG VERSION

ENV HOME=/opt/helm
COPY watches.yaml ${HOME}/watches.yaml
COPY helm-charts  ${HOME}/helm-charts
WORKDIR ${HOME}

COPY LICENSE /licenses/

LABEL name="NGINX Ingress Operator" \
      maintainer="kubernetes@nginx.com" \
      vendor="NGINX Inc" \
      version="v${VERSION}" \
      release="1" \
      summary="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      description="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers"
