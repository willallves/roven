.ONESHELL:

EKS_DIR ?= ./dev/kubernetes
BINARIES ?= hybrid_server hybrid_agent
OSES ?= linux
ARCHITECTURES ?= amd64
VERSION ?= latest
DOCKER_HUB ?= public.ecr.aws/n9c4h4j5
DOCKER_TAG_AGENT ?= hybrid-attestor-server:latest hybrid-attestor-agent:latest 
BUILD_DIR ?= ./build
PLATFORMS ?= $(foreach os, $(OSES), $(foreach architecture, $(ARCHITECTURES), --platform $(os)/$(architecture)))
REPOSITORY_SERVER ?= hybrid-attestor-server 
REPOSITORY_AGENT ?= hybrid-attestor-agent
IMAGE_TAG ?= latest

BUILD_TARGETS := $(foreach binary, $(BINARIES), $(foreach os, $(OSES), $(foreach architecture, $(ARCHITECTURES), $(binary)-$(os)-$(architecture))))
DOCKER_TARGETS := $(foreach binary, $(BINARIES), $(binary)-docker)

target_words = $(subst -, ,$@)
target_binary = $(word 1, $(target_words))
target_os = $(word 2, $(target_words))
target_architecture = $(word 3, $(target_words))
target_software_type = $(word 2, $(subst _, ,$(target_binary)))

target_binary_hyphens = $(subst _,-attestor-,$(target_binary))

build: $(BUILD_TARGETS)
$(BUILD_TARGETS):
	CGO_ENABLED=0 GOOS=$(target_os) GOARCH=$(target_architecture) go build -ldflags="-s -w -extldflags -static" -o $(BUILD_DIR)/$(target_os)/$(target_architecture)/$(target_binary) cmd/$(target_binary)/main.go

test: test-unit test-integration
	go test ./...

test-unit:
	go test ./...

test-integration:
	bash ./test/integration/test.sh

docker:
	docker build -f hybrid/dev/docker/server.Dockerfile -t $DOCKER_HUB/$REPOSITORY_SERVER:$IMAGE_TAG .
	docker build -f hybrid/dev/docker/agent.Dockerfile -t $DOCKER_HUB/$REPOSITORY_AGENT:$IMAGE_TAG .
	docker push $DOCKER_HUB/$REPOSITORY_SERVER:$IMAGE_TAG
	docker push $DOCKER_HUB/$REPOSITORY_AGENT:$IMAGE_TAG

docker-build:
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -ldflags="-s -w -extldflags -static" -o ${BINARY} cmd/${BINARY}/main.go

deploy-spire-eks:
	kubectl delete --all daemonsets.app --namespace=spire
	kubectl delete --all statefulset.app --namespace=spire
	kubectl delete --ignore-not-found namespace spire
	envsubst < $(EKS_DIR)/server.yaml | kubectl apply -f -
	envsubst < $(EKS_DIR)/agent.yaml | kubectl apply -f -

clean:
	rm -rf $(BUILD_DIR)

.PHONY: $(BUILD_TARGETS) $(DOCKER_TARGETS) build test docker clean
