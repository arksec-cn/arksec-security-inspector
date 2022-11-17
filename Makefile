VERSION?=3.0.0
BUILD_DATE := $(shell date +"%Y-%m-%d_%H:%I:%S")
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	LDFLAGS = -extldflags "-static"
endif
OS_NAME?=linux
OS_ARCH?=amd64
CURPATH?=$(PWD)
SOURCEPATH=$(PWD)
RUNTIME?=docker
BUILD_TAG?=
PROJECT_ENV?=dev
PROJECT_NAME?="virgo"
IMAGE_TAG?=latest

IMAGE=registry.arksec.cn/virgo/$(PROJECT_ENV)/$(OS_ARCH)/inspector:${IMAGE_TAG}

.PHONY:inspector
inspector:
	@echo "+ $@"
	CGO_ENABLED=0 GOOS=$(OS_NAME) GOARCH=$(OS_ARCH) GOPROXY="https://goproxy.cn,direct" go build -mod=mod \
 		-o dist/inspector cmd/inspector/main.go

.PHONY:inspector-build
inspector-build:
	@echo "+ $@"
	$(RUNTIME) build -f build/inspector/Dockerfile --build-arg PROJECT_ENV=$(PROJECT_ENV) \
	--build-arg PROJECT_NAME=$(PROJECT_NAME) \
	--build-arg OS_ARCH=$(OS_ARCH) \
	-t $(IMAGE) .

.PHONY:inspector-push
inspector-push:
	@echo "+ $@"
	$(RUNTIME) push $(IMAGE)

push-inspector:
	make inspector
	make inspector-build
	make inspector-push