.DEFAULT: build-binary
.PHONY: build-binary build-image clean

export PATH:=$(PATH):$(CURDIR)"/build/_output/bin"

REGISTRY?=docker.io
ORG?=projectraksh

DOCKERFILES=Dockerfile.sc-scratch Dockerfile.securecontainer-operator
IMAGE?=securecontainer-operator
IMAGES=$(subst Dockerfile.,,$(DOCKERFILES))
IMAGE_NAME=$(ORG)/$(IMAGE)

VERSION?=master
BUILD_FLAGS?=-a -tags netgo -ldflags '-w -extldflags "-static"'
DOCKER ?= docker
SED_I?=sed -i
GOHOSTOS ?= $(shell go env GOHOSTOS)

ifeq ($(GOHOSTOS),darwin)
  SED_I=sed -i ''
endif

ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)
ALL_ARCH = amd64 ppc64le

QEMUVERSION=v3.0.0

BASEIMAGE=alpine:3.9

ifeq ($(ARCH),x86_64)
	override ARCH=amd64
endif

ifeq ($(ARCH),ppc64le)
    QEMUARCH=ppc64le
	override BASEIMAGE=ppc64le/alpine:3.9
endif

TEMP_DIR := $(shell mktemp -d)

register:
	docker run --rm --privileged multiarch/qemu-user-static:register || true

build-image: $(addprefix build-image-,$(IMAGES))

build-image-%:
	$(MAKE) IMAGE=$* all-container

push-image: $(addprefix push-image-,$(IMAGES))

push-image-%:
	$(MAKE) IMAGE=$* all-push

all-container: register $(addprefix sub-container-,$(ALL_ARCH))

sub-container-%:
	$(MAKE) ARCH=$* container

container: .container-$(ARCH)
.container-$(ARCH): build-binary-$(ARCH)
	cp -r ./build $(TEMP_DIR)
	cd $(TEMP_DIR) && $(SED_I) "s|BASEIMAGE|$(BASEIMAGE)|g" build/Dockerfile.$(IMAGE)
	cd $(TEMP_DIR) && $(SED_I) "s|ARCH|$(QEMUARCH)|g" build/Dockerfile.$(IMAGE)

ifeq ($(ARCH),amd64)
	# When building "normally" for amd64, remove the whole line, it has no part in the amd64 image
	cd $(TEMP_DIR) && $(SED_I) "/CROSS_BUILD_/d" build/Dockerfile.$(IMAGE)
else
	# When cross-building, only the placeholder "CROSS_BUILD_" should be removed
	curl -sSL https://github.com/multiarch/qemu-user-static/releases/download/$(QEMUVERSION)/x86_64_qemu-$(QEMUARCH)-static.tar.gz | tar -xz -C $(TEMP_DIR)
	cd $(TEMP_DIR) && $(SED_I) "s/CROSS_BUILD_//g" build/Dockerfile.$(IMAGE)
endif
	$(DOCKER) build --no-cache -t $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH) -f $(TEMP_DIR)/build/Dockerfile.$(IMAGE) $(TEMP_DIR)
	$(DOCKER) tag $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH) $(REGISTRY)/$(IMAGE_NAME):latest-$(ARCH)

build-binary-%: clean
	$(MAKE) OS=linux build-binary

build-binary:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILD_FLAGS) -o build/_output/bin/securecontainer-operator cmd/manager/main.go
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILD_FLAGS) -o build/_output/bin/rakshctl cmd/rakshctl/main.go
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILD_FLAGS) -o build/_output/bin/sc-scratch cmd/sc-scratch/main.go

tests: unit functional

unit:
	go test -v ./pkg/...

functional: build-binary
	PATH=$(PATH) go test -v ./...

all-push: $(addprefix sub-push-,$(ALL_ARCH))

sub-push-%:
	$(MAKE) ARCH=$* push

docker-login:
	docker login $(REGISTRY) -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

push:
	$(DOCKER) push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH)
	$(DOCKER) push $(REGISTRY)/$(IMAGE_NAME):latest-$(ARCH)

push-manifest:
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create --amend $(REGISTRY)/$(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-amd64 $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-ppc64le
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate --arch amd64 $(REGISTRY)/$(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-amd64
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate --arch ppc64le $(REGISTRY)/$(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)-ppc64le
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push --purge $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

verify:
	./hack/verify-all.sh

clean:
	rm -rf build/_output
