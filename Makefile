TAG = latest
IMAGE = nginx-ingress-operator
BUILD_FLAGS =

RUN_NAMESPACE = default

test:
	GO111MODULE=on go test ./...

binary:
	CGO_ENABLED=0 GO111MODULE=on $(BUILD_FLAGS) GOOS=linux go build -installsuffix cgo -o build/_output/bin/nginx-ingress-operator github.com/nginxinc/nginx-ingress-operator/cmd/manager

build: binary
	docker build -f build/Dockerfile -t $(IMAGE):$(TAG) .

run-local:
	operator-sdk run --local --namespace=$(RUN_NAMESPACE)

generate-crds:
	operator-sdk generate k8s && operator-sdk generate crds

lint:
	golangci-lint run

.PHONY: build