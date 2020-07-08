OLD_TAG = 0.0.3
TAG = 0.0.4

IMAGE = nginx-ingress-operator

test:
	GO111MODULE=on go test ./...

binary:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -trimpath -installsuffix cgo -o build/_output/bin/nginx-ingress-operator github.com/nginxinc/nginx-ingress-operator/cmd/manager

build: binary
	docker build -f build/Dockerfile -t $(IMAGE):$(TAG) .

run-local:
	operator-sdk run local

generate-crds:
	operator-sdk generate k8s && operator-sdk generate crds

lint:
	golangci-lint run

generate-metadata: generate-crds
	operator-sdk generate csv --csv-version=$(TAG) --from-version=$(OLD_TAG) --make-manifests=false
	echo "Metadata generated, please make sure you add/update fields in nginx-ingress-operator.v$(TAG).clusterserviceversion.yaml"

generate-bundle:
	-rm -rf bundle
	mkdir bundle
	cp deploy/crds/* bundle/
	cp deploy/olm-catalog/nginx-ingress-operator/nginx-ingress-operator.package.yaml bundle/
	./hack/copy_manifests.sh
	-rm bundle.zip
	zip -j bundle.zip bundle/*.yaml

.PHONY: build