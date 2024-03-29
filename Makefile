SHELL := /bin/bash -o pipefail

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

TMP_BASE := .tmp
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)/bin
TMP_VERSIONS := $(TMP)/versions

export GO111MODULE := on
export GOBIN := $(abspath $(TMP_BIN))
export PATH := $(GOBIN):$(PATH)

# This is the only variable that ever should change.
# This can be a branch, tag, or commit.
# When changed, the given version of Prototool will be installed to
# .tmp/$(uname -s)/(uname -m)/bin/prototool
PROTOTOOL_VERSION := v1.8.0

PROTOTOOL := $(TMP_VERSIONS)/prototool/$(PROTOTOOL_VERSION)
$(PROTOTOOL):
	$(eval PROTOTOOL_TMP := $(shell mktemp -d))
	cd $(PROTOTOOL_TMP); go get github.com/uber/prototool/cmd/prototool@$(PROTOTOOL_VERSION)
	@rm -rf $(PROTOTOOL_TMP)
	@rm -rf $(dir $(PROTOTOOL))
	@mkdir -p $(dir $(PROTOTOOL))
	@touch $(PROTOTOOL)

# proto is a target that uses prototool.
# By depending on $(PROTOTOOL), prototool will be installed on the Makefile's path.
# Since the path above has the temporary GOBIN at the front, this will use the
# locally installed prototool.
.PHONY: proto

generate: $(PROTOTOOL)
	@go get github.com/grpc-ecosystem/grpc-gateway@v1.11.3
	@go install ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/protoc-gen-grpc-gateway
	@go install ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/protoc-gen-swagger
	@go get github.com/fiorix/protoc-gen-cobra@v0.0.0-20181029091941-dffa0bfa45cc
	@prototool lint
	@prototool generate
	@rm -r .tmp

doc:
	@npm install -g redoc-cli
	@redoc-cli bundle \
		./openapi/vinyltap.swagger.json \
		-o="./openapi/index.html" \
		--title "Vinyl Tap API"

# do not include `generate` in the docker command, as the Dockerfile
# runs code generation on build.
docker:
	@docker build --rm -t rpc-demo .

start: generate doc
	@go run ./cmd/api/main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
