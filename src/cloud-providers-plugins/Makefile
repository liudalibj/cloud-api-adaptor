# (C) Copyright Confidential Containers Contributors
# SPDX-License-Identifier: Apache-2.0

.PHONY: all build libvirt-ext image-plugins image-cca image-ppc images

SHELL = bash -o pipefail

ARCH        ?= $(subst x86_64,amd64,$(shell uname -m))
ARCHES      ?= linux/amd64,linux/s390x
GOOPTIONS   ?= GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=1
GOFLAGS     ?=
SOURCEDIRS  := ./libvirt-ext
SOURCES     := $(shell find $(SOURCEDIRS) -name '*.go' -print)

BUILTIN_CLOUD_PROVIDERS ?= libvirt


.PHONY: .git-commit.tmp
.git-commit: .git-commit.tmp
	@cmp $< $@ >/dev/null 2>&1 || cp $< $@
.git-commit.tmp:
	@printf "$$(git rev-parse HEAD 2>/dev/null || echo unknown)" >$@
	@test -n "$$(git status --porcelain --untracked-files=no 2> /dev/null)" && echo -dirty >>$@ || true

REGISTRY ?= quay.io/confidential-containers
DEV_TAGS ?= $(shell cat .git-commit)
PLUGINS_IMAGE ?= $(REGISTRY)/cloud-provider-plugins:$(DEV_TAGS)
PLUGINS_CAA_IMAGE ?= $(REGISTRY)/cloud-api-adaptor:dev-plugins-$(DEV_TAGS)
PLUGINS_PEERPOD_CTRL_IMAGE ?= $(REGISTRY)/peerpod-ctrl:dev-plugins-$(DEV_TAGS)
CAA_IMG ?= quay.io/confidential-containers/cloud-api-adaptor:latest
PEERPOD_CTRL_IMG ?= quay.io/confidential-containers/peerpod-ctrl:latest
PUSH ?= true
# If not pushing `--load` into the local docker cache
DOCKER_OPTS := $(if $(filter $(PUSH),true),--push,--load)


all: build
build: libvirt-ext

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

space := $() $()
comma := ,
GOFLAGS += -tags=$(subst $(space),$(comma),$(strip $(BUILTIN_CLOUD_PROVIDERS)))

libvirt-ext: $(SOURCES)
	$(GOOPTIONS) go build $(GOFLAGS) -buildmode=plugin ./$@

.PHONY: image-plugins
image-plugins: .git-commit ## Build and push cloud provider plugins
	cd ../ && docker buildx build --platform "$(ARCHES)" \
	-t $(PLUGINS_IMAGE) \
	-f cloud-providers-plugins/Dockerfile.plugins \
	$(DOCKER_OPTS) .

.PHONY: image-caa
image-caa: .git-commit ## Build and push caa with cloud provider plugins
	cd ../ && docker buildx build --platform "$(ARCHES)" \
	-t $(PLUGINS_CAA_IMAGE) \
	-f cloud-providers-plugins/Dockerfile.caa \
	--build-arg PLUGINS_IMG=$(PLUGINS_IMAGE) \
	--build-arg CAA_IMG=$(CAA_IMG) \
	$(DOCKER_OPTS) .

.PHONY: image-ppc
image-ppc: .git-commit ## Build and push peerpod-ctrl with cloud provider plugins
	cd ../ && docker buildx build --platform "$(ARCHES)" \
	-t $(PLUGINS_PEERPOD_CTRL_IMAGE) \
	-f cloud-providers-plugins/Dockerfile.peerpod_ctrl \
	--build-arg PLUGINS_IMG=$(PLUGINS_IMAGE) \
	--build-arg PEERPOD_CTRL_IMG=$(PEERPOD_CTRL_IMG) \
	$(DOCKER_OPTS) .

.PHONY: images
images: image-plugins image-caa image-ppc
