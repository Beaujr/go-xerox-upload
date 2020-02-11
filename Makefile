PACKAGE_NAME := github.com/beaujr/go-xerox-upload
REGISTRY := docker.io
APP_NAME := beaujr/go-xerox-upload
IMAGE_TAG ?= 0.1
GOPATH ?= $HOME/go
BUILD_TAG := build
BINPATH := ./bin
NAMESPACE := default

# Path to dockerfiles directory
DOCKERFILES := build

# Go build flags
GOOS := linux
GOARCH := amd64
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_SHORT_COMMIT := $(shell git rev-parse --short HEAD)
GOLDFLAGS := -ldflags "-X $(PACKAGE_NAME)/pkg/util.AppGitCommit=${GIT_COMMIT} -X $(PACKAGE_NAME)/pkg/util.AppVersion=${IMAGE_TAG}"

.PHONY: verify build docker_build push generate generate_verify \
	go_upload_xerox go_test go_fmt e2e_test go_verify   \
	docker_build docker_push

# Alias targets
###############

build: go_dep go_test go_upload_xerox # docker_build
verify: generate_verify go_verify
#push: build docker_push

# Go targets
#################
go_verify: go_fmt go_test

go_dep:
	dep ensure -v

go_upload_xerox:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-a -tags netgo \
		-o $(BINPATH)/${APP_NAME}-$(GOOS)_$(GOARCH) \
		./

go_test:
ifeq ($(GOARCH),amd64)
	CGO_ENABLED=0 go test -v \
		-cover \
		-coverprofile=coverage.out \
		$$(go list ./... | \
			grep -v '/vendor/' | \
			grep -v '/pkg/client' \
		)
endif

coverage: go_test
	go tool cover -html=coverage.out

go_fmt:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi


# Docker targets
################
docker_build:
	docker build \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOOS=$(GOOS) \
		-t $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) \
		-f $(DOCKERFILES)/Dockerfile \
		./

docker_run:
	@docker run -v $(shell pwd)/files:/tmp -p 8080:10000 $(REGISTRY)/$(APP_NAME):$(BUILD_TAG)

docker_push: docker-login
	set -e; \
	docker tag $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) $(APP_NAME):$(IMAGE_TAG)-$(GOARCH)-$(GIT_SHORT_COMMIT) ; \
	docker push $(APP_NAME):$(IMAGE_TAG)-$(GOARCH)-$(GIT_SHORT_COMMIT);
ifeq ($(GITHUB_HEAD_REF),master)
	docker tag $(APP_NAME):$(IMAGE_TAG)-$(GOARCH)-$(GIT_SHORT_COMMIT) $(APP_NAME):latest_$(GOARCH)
	docker push $(APP_NAME):latest_$(GOARCH)
endif

check-docker-credentials:
ifndef DOCKER_USER
	$(error DOCKER_USER is undefined)
else
  ifndef DOCKER_PASS
	$(error DOCKER_PASS is undefined)
  endif
endif

docker-login: check-docker-credentials
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS) $(REGISTRY)

score: PR_ID=$(shell echo $(GITHUB_REF) | tr -dc '0-9')
score:
	curl -X GET \
	https://gogitops.beau.cf/$(GITHUB_REPOSITORY)/pull/$(PR_ID) \
	-H 'apikey: $(GITOPS_API_KEY)' \
	-H 'user: $(GITHUB_USER)' \
	-H 'token: $(GITHUB_TOKEN)'

deploy:
	docker build -t gcloud -f build/Dockerfile.deploy .; \
	docker run -e GCLOUD_API_KEYFILE=$(GCLOUD_API_KEYFILE) \
	-e CLOUDSDK_CORE_PROJECT=go-xerox-upload \
	-e CLIENT_ID=$(CLIENT_ID) \
	-e PROJECT_ID=$(PROJECT_ID) \
	-e CLIENT_SECRET=$(CLIENT_SECRET) \
	-e ACCESS_TOKEN=$(ACCESS_TOKEN) \
	-e REFRESH_TOKEN=$(REFRESH_TOKEN) \
	-e EXPIRY_TIME=$(EXPIRY_TIME) gcloud:latest;
