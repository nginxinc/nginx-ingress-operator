#!/bin/bash

image=$1
version=$2

kube_image=kubebuilder/kube-rbac-proxy
kube_image_version=v0.8.0

token="$(curl 'https://auth.docker.io/token?service=registry.docker.io&scope=repository:'${image}':pull' 2>/dev/null | jq -r '.token')"

image_digest=$(curl -sSfL -I -H "Authorization: Bearer ${token}" -H "Accept: application/vnd.docker.distribution.manifest.list.v2+json" "https://index.docker.io/v2/${image}/manifests/${version}" | awk 'BEGIN {FS=": "}/^docker-content-digest/{gsub(/"/, "", $2); print $2}')

digest="$(curl -sSfL -H "Authorization: Bearer ${token}" -H "Accept: application/vnd.docker.distribution.manifest.v2+json" "https://index.docker.io/v2/${image}/manifests/${version}" | jq -r '.config.digest')"

created=$(curl -sSfL -H "Accept: application/vnd.docker.distribution.manifest.v2+json" -H "Authorization: Bearer ${token}" "https://index.docker.io/v2/${image}/blobs/${digest}" | jq -r '.config.Labels."org.opencontainers.image.created"')

proxy="./config/default/manager_auth_proxy_patch.yaml"
kube_proxy=$(yq e '.spec.template.spec.containers.[0].image' $proxy)
full_image=${kube_proxy%:*}
kube_image=${full_image#*/}
kube_version=${kube_proxy#*:}

kube_digest=$(curl -sSfL -I -H "Accept: application/vnd.docker.distribution.manifest.list.v2+json" "https://gcr.io/v2/${kube_image}/manifests/${kube_version}" | awk 'BEGIN {FS=": "}/^docker-content-digest/{gsub(/"/, "", $2); print $2}')

printf "%s\n\n" "Manually repleace the following values in bundle/manifests/nginx-ingress-operator.clusterserviceversion.yaml"
printf "%s\n" "metadata.annotations.createdAt: ${created}"
printf "%s\n" "metadata.annotations.containerImage: docker.io/${image}@${image_digest}"
printf "%s\n" ".spec.install.spec.deployments[0].spec.template.spec.containers[1].image (nginx-ingress-operator): docker.io/${image}@${image_digest}"
printf "%s\n" ".spec.install.spec.deployments[0].spec.template.spec.containers[0].image (kube-rbac-proxy): ${full_image}@${kube_digest}"
