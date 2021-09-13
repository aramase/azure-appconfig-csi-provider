REPOSITORY ?= aramase/azure-appconfig-csi-provider
IMG := $(REPOSITORY):latest
ARCH ?= "linux/amd64"

docker-build:
	docker buildx build --platform=${ARCH} -t ${IMG} . --load

docker-build-push:
	docker buildx build --platform=${ARCH} -t ${IMG} . --push
