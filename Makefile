OLD_TAG = 0.1.0
TAG = 0.2.0

IMAGE = nginx/nginx-ingress-operator

test:
	GO111MODULE=on go test ./...

binary:
	$(eval GOPATH=$(shell go env GOPATH))
	CGO_ENABLED=0 GO111MODULE=on GOFLAGS="-gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH}" GOOS=linux go build -trimpath -ldflags "-s -w" -o build/_output/bin/nginx-ingress-operator github.com/nginxinc/nginx-ingress-operator/cmd/manager

build: binary
	docker build -f build/Dockerfile -t $(IMAGE):$(TAG) .

run-local:
	go run github.com/operator-framework/operator-sdk/cmd/operator-sdk run local

generate-crds:
	go run github.com/operator-framework/operator-sdk/cmd/operator-sdk generate k8s && operator-sdk generate crds --crd-version v1beta1

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

generate-metadata: generate-crds
	go run github.com/operator-framework/operator-sdk/cmd/operator-sdk generate csv --csv-version=$(TAG) --from-version=$(OLD_TAG) --make-manifests=false
	echo "Metadata generated, please make sure you add/update fields in nginx-ingress-operator.v$(TAG).clusterserviceversion.yaml"

generate-bundle:
	@mkdir bundle/$(TAG)
	@cp deploy/crds/* bundle/$(TAG)
	@cp deploy/olm-catalog/nginx-ingress-operator/nginx-ingress-operator.package.yaml bundle/
	@cp -R deploy/olm-catalog/nginx-ingress-operator/$(TAG)/ bundle/$(TAG)/
	cd bundle && opm alpha bundle generate -d ./$(TAG)/ -u ./$(TAG)/
	@mv bundle/bundle.Dockerfile bundle/bundle-$(TAG).Dockerfile
	@echo '\nLABEL com.redhat.openshift.versions="v4.5,v4.6"\nLABEL com.redhat.delivery.operator.bundle=true\nLABEL com.redhat.delivery.backport=true' >> bundle/bundle-$(TAG).Dockerfile
	docker build -t bundle:$(TAG) -f bundle/bundle-$(TAG).Dockerfile ./bundle

.PHONY: build
