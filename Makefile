PACKAGE_NAME := github.com/beaujr/go-xerox-upload
REGISTRY := docker.com
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
GOLDFLAGS := -ldflags "-X $(PACKAGE_NAME)/pkg/util.AppGitCommit=${GIT_COMMIT} -X $(PACKAGE_NAME)/pkg/util.AppVersion=${IMAGE_TAG}"

.PHONY: verify build docker_build push generate generate_verify \
	go_upload_xerox go_test go_fmt e2e_test go_verify   \
	docker_build docker_push

# Alias targets
###############

build: go_dep  go_test go_upload_xerox # docker_build
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
	go test -v \
		-cover \
		-coverprofile=coverage.out \
		$$(go list ./... | \
			grep -v '/vendor/' | \
			grep -v '/pkg/client' \
		)

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
		-t $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) \
		-f $(DOCKERFILES)/Dockerfile \
		./

#docker_push:
#	set -e; \
#	docker tag $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) $(REGISTRY)/$(APP_NAME):$(IMAGE_TAG) ; \
#	docker push $(REGISTRY)/$(APP_NAME):$(IMAGE_TAG);
#