.PHONY: test clean docker-builder docker-test docker-build setup help

export PATH := ./bin:$(PATH)

prototool:
	@type prototool >/dev/null 2>&1 || GO111MODULE=off go get github.com/uber/prototool/cmd/prototool@v1.7.0

generate: prototool
	@prototool format -f -w
	@prototool lint
	@prototool generate

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help