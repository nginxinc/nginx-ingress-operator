TAG = latest
IMAGE = nginx-ingress-operator

RUN_NAMESPACE = default

test:
	GO111MODULE=on go test ./...

build:
	operator-sdk build $(IMAGE):$(TAG)

run-local:
	operator-sdk run --local --namespace=$(RUN_NAMESPACE)

generate-crds:
	operator-sdk generate k8s && operator-sdk generate crds

lint:
	golangci-lint run

.PHONY: build