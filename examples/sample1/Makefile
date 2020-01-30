PROJECTNAME := sample1
VERSION := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILD   := $(shell git rev-parse --short HEAD)
MODULE  := $(shell basename "$(PWD)")
BUILD_DATE=$(shell date +%Y-%m-%d_%H:%M)

##############################################################################################

# Setup the -ldflags option for go build here, interpolate the variable values
# -w : omit the DWARF symbol table during the build process
# -s : omit the symbol table and debugging information
LDFLAGS = -ldflags "-w -s -X main.Version=${VERSION}:${BUILD} -X main.Build=${BUILD_DATE} -X main.Module=${MODULE}"

##############################################################################################

# overwritable with inline env variable
TAG ?= latest
IMAGE ?= $(PROJECTNAME)
IMAGE_REGISTRY ?= techlab21.azurecr.io

IMAGE_FULL_PATH=$(IMAGE_REGISTRY)/$(IMAGE)

##############################################################################################

# Enable go modules
#GOCMD = GO111MODULE=on go
GOCMD = go

##############################################################################################
.PHONY: build
build:   ## build the binary file for server
	${GOCMD} build ${LDFLAGS} -o ${PROJECTNAME}

build-app:    ## build the binary in build/app
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GOCMD} build ${LDFLAGS} -o build/app

dep: ## Get the dependencies
	@go get -v -d ./...

.PHONY: run-local
run-local:
	go run main.go

##############################################################################################
build-image: build-app    ## build container image
	docker build -t $(IMAGE_FULL_PATH):$(TAG) .

test-image:  build-image   ## build image and launch local image for testing
    # --rm to automatically delete it when it exits
	docker run --rm --name $(IMAGE)_test -p 9000:80 --env PORT=80 $(IMAGE_FULL_PATH):$(TAG)

push-image:  build-image   ## push container image to registry $(IMAGE_REGISTRY)
	docker push $(IMAGE_FULL_PATH):$(TAG)

##############################################################################################

deploy: push-image    ## ENV INTEGRATION
	@kubectl -n $(PROJECTNAME)-int scale deployment website --replicas=0
	@kubectl -n $(PROJECTNAME)-int scale deployment website --replicas=1

deploy-hml: push-image  ## ENV HOMOLOGATION
	@kubectl -n $(PROJECTNAME)-hml scale deployment website --replicas=0
	@kubectl -n $(PROJECTNAME)-hml scale deployment website --replicas=1

deploy-prd:    ## ENV PRODUCTION
	@kubectl -n $(PROJECTNAME) scale deployment website --replicas=0
	@kubectl -n $(PROJECTNAME) scale deployment website --replicas=1

##############################################################################################

publish-init:   ## init K8s env
	@kubectl apply -f ./kube/config_k8s-namespace.yaml
	@kubectl -n $(PROJECTNAME) create secret generic imageregistry --from-file=.dockerconfigjson=./kube/config_auth.json --type=kubernetes.io/dockerconfigjson | true
	@kubectl apply -f ./kube/config_k8s.yaml

## CREATE SECRET ON ERROR
## kubectl create secret generic imageregistry --from-file=.dockerconfigjson=$HOME/.docker/config.json  --type=kubernetes.io/dockerconfigjson --namespace $(PROJECTNAME)

.PHONY: help
help: ## Display this help screen
	@echo " VERSION:" $(VERSION) " BUILD:" $(BUILD) " PROJECTNAME:" $(PROJECTNAME) " MODULE:" $(MODULE)
	@echo " IMAGE:" $(IMAGE_FULL_PATH):$(TAG)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

