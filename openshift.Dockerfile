# Build the manager binary
FROM golang:1.17.4 as builder
ARG VERSION

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=${VERSION}" -a -o manager main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
COPY --from=builder /workspace/manager .
COPY config/crd/kic ./config/crd/kic
COPY LICENSE /licenses/

LABEL name="NGINX Ingress Operator" \
      description="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      summary="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      io.openshift.tags="nginx,ingress-controller,ingress,controller,kubernetes,openshift" \
      maintainer="NGINX Inc <kubernetes@nginx.com>" \
      vendor="NGINX Inc <kubernetes@nginx.com>"

ENTRYPOINT ["/manager"]
