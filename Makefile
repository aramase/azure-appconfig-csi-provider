REPOSITORY ?= aramase/azure-appconfig-csi-provider
IMG := $(REPOSITORY):latest
ARCH ?= "linux/amd64"

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o _output/azure-appconfig-csi-provider main.go

docker-build: clean build
	docker buildx build --no-cache --platform=${ARCH} -t ${IMG} -f Dockerfile . --load

docker-build-push:
	docker buildx build --no-cache --platform=${ARCH} -t ${IMG} -f Dockerfile . --push

clean:
	rm -rf _output
