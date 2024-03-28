# (C) Copyright Confidential Containers Contributors
# SPDX-License-Identifier: Apache-2.0

.PHONY: all build libvirt-ext image-plugins image-cca image-ppc images

SHELL = bash -o pipefail

ARCH        ?= $(subst x86_64,amd64,$(shell uname -m))

GOOPTIONS   ?= GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=1
GOFLAGS     ?=
SOURCEDIRS  := ./libvirt-ext
SOURCES     := $(shell find $(SOURCEDIRS) -name '*.go' -print)

BUILTIN_CLOUD_PROVIDERS ?= libvirt

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
image-plugins: ## Build and push cloud provider plugins
	./build.sh -p

.PHONY: image-caa
image-caa: ## Build and push caa with cloud provider plugins
	./build.sh -c

.PHONY: image-ppc
image-ppc: ## Build and push peerpod-ctrl with cloud provider plugins
	./build.sh -l

.PHONY: images
images: image-plugins image-cca image-ppc
