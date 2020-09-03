# Build the manager binary
FROM golang:1.13 as builder

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go


FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
LABEL name="NGINX Ingress Operator" \
      description="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      summary="The NGINX Ingress Operator is a Kubernetes/OpenShift component which deploys and manages one or more NGINX/NGINX Plus Ingress Controllers" \
      io.openshift.tags="nginx,ingress-controller,ingress,controller,kubernetes,openshift" \
      maintainer="NGINX Inc <kubernetes@nginx.com>" \
      vendor="NGINX Inc <kubernetes@nginx.com>"
WORKDIR /
COPY --from=builder /workspace/manager .
COPY LICENSE /licenses/
COPY hack/kic_crds /kic_crds
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
