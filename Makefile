export CGO_ENABLED=0

REGION?=ap-northeast-2
REGISTRY_ID?=246092231098
IMAGE_NAME?=dax-custom-webhook
IMAGE?=$(REGISTRY_ID).dkr.ecr.$(REGION).amazonaws.com/$(IMAGE_NAME)
TAG?=$(shell git describe --tags $(shell git rev-list --tags --max-count=1))

.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/dax-custom-webhook-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/dax-custom-webhook-amd64 .

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: docker-build
docker-build:
	@echo '> Building docker image $(IMAGE)'
	git remote update
	docker build --no-cache -t $(IMAGE):$(TAG) .
	docker tag  $(IMAGE):$(TAG) $(IMAGE):latest

.PHONY: docker-push
docker-push:
	@echo '> Pushing docker image $(IMAGE)'
	if ! aws ecr get-login-password --region $(REGION) | docker login --username AWS --password-stdin $(REGISTRY_ID).dkr.ecr.$(REGION).amazonaws.com; then \
	  eval $$(aws ecr get-login --registry-ids $(REGISTRY_ID) --no-include-email); \
	fi
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

.PHONY: docker
docker: docker-build docker-push

# `kind` is required
.PHONY: create-cluster
create-cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster

.PHONY: load-image
load-image:
	@echo "\n🔧 Loading $(IMAGE)..."
	kind load docker-image $(IMAGE):latest 

# `kind` is required
.PHONY: delete-cluster
delete-cluster:
	@echo "\n🔧 Deleting Kubernetes cluster..."
	kind delete cluster

.PHONY: init
init:
	@echo '> Creating signed cert'
	./hack/webhook-gen-certs.sh 

.PHONY: deploy-webhook
deploy:
	kubectl apply -f ./hack/manifests/webhook

.PHONY: delete-webhook
delete:
	kubectl delete -f ./hack/manifests/webhook

.PHONY: logs
logs:
	@echo "\n🔍 Streaming $(IMAGE) logs..."
	kubectl logs -l app=$(IMAGE) -n kube-system -f

