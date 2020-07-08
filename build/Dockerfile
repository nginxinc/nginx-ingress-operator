FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/nginx-ingress-operator \
    USER_UID=1001 \
    USER_NAME=nginx-ingress-operator

# install operator binary
COPY build/_output/bin/nginx-ingress-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

COPY LICENSE /licenses/
COPY build/kic_crds /kic_crds

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
LABEL name="NGINX Ingress Operator" \
      description="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      summary="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      io.openshift.tags="nginx,ingress-controller,ingress,controller,kubernetes,openshift" \
      maintainer="NGINX Inc <kubernetes@nginx.com>" \
      vendor="NGINX Inc <kubernetes@nginx.com>"