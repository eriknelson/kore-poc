REGISTRY         ?= quay.io
ORG              ?= hegemone
TAG              ?= latest
IMAGE            ?= $(REGISTRY)/$(ORG)/korecomm-go-poc
BUILD_DIR        = "${GOPATH}/src/github.com/hegemone/kore-poc/korecomm-go/build"
SOURCES          := $(shell find . -name '*.go' -not -path "*/vendor/*" -not -path "*/extensions/*")
#PLUGIN_SOURCES   := $(shell find . -name '*/extensions/plugin/*.go' -not -path "*/vendor/*")
.DEFAULT_GOAL    := build
PROJECT_ROOT := ""$(abspath $(lastword $(MAKEFILE_LIST)))/..""

vendor:
	@glide install -v

plugins: $(PLUGIN_SOURCES)
	@go build -buildmode=plugin -o ${BUILD_DIR}/bacon.plugins.kore.nsk.io.so -i -ldflags="-s -w" ./pkg/extension/plugin/bacon.go

korecomm: $(SOURCES) plugins
	@go build -o ${BUILD_DIR}/korecomm -i -ldflags="-s -w" ./cmd/korecomm

build: korecomm
	@echo > /dev/null

clean:
	@rm -rf ${BUILD_DIR}

run: korecomm
	@KORECOMM_PLUGIN_DIR=${PROJECT_ROOT}/build \
	KORECOMM_ADAPTER_DIR=${PROJECT_ROOT}/build \
	./build/korecomm

image:
	docker build -t ${IMAGE}:${TAG} ${PROJECT_ROOT}

run-image:
	docker run -it ${IMAGE}:${TAG}

push:
	docker push ${IMAGE}:${TAG}

#test:
	#go test ./pkg/...

.PHONY: vendor image push clean run image run-image push
