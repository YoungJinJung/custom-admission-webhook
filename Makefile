.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/admission-webhook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building Docker image: custom-admission-webhook"
	docker build --no-cache -t custom-admission-webhook:latest .

.PHONY: clean
clean:
	rm -rf ./bin

# `kind` is required
.PHONY: create-cluster
create-cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster

.PHONY: load-image
load-image:
	@echo "\n🔧 Loading custom-admission-webhook..."
	kind load docker-image custom-admission-webhook:latest 

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
	@echo "\n🔍 Streaming custom-admission-webhook logs..."
	kubectl logs -l app=custom-admission-webhook -n kube-system -f

