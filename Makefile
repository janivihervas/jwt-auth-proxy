SHELL := /bin/bash

VERSION ?= dev
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
TAG ?= $(shell git describe --tags 2>/dev/null)

# If VERSION is not it's default value dev, change it for corresponding git id.
# Will not have effect if VERSION is overwritten with VERSION=<version> make [target...]
ifeq ($(VERSION), dev)
  ifneq ($(TAG),)
  	VERSION = $(TAG)
  else ifeq ($(BRANCH), master)
    VERSION = latest
  endif
endif

# Maximum parallel jobs. Defaults to the number of physical cores on the machine.
# Override with
#   PARALLELISM=1 make [target...]
PARALLELISM ?= 0
ifeq ($(PARALLELISM), 0)
  ifeq ($(shell uname), Darwin)
    PARALLELISM = $(shell sysctl hw.physicalcpu | tr " " "\n" | tail -n 1)
  else ifeq ($(shell uname), Linux)
    PARALLELISM = $(shell cat /proc/cpuinfo | grep "model name" | wc -l)
  else
    PARALLELISM = 1
  endif
endif

MAKEFLAGS += --no-print-directory --output-sync --jobs=$(PARALLELISM)

PACKAGE := github.com/janivihervas/authproxy

GITHUB_API_URL := https://api.github.com/repos/janivihervas/authproxy

# These are all the combinations of OS and architecture we're compiling to. See
# https://golang.org/doc/install/source#environment for full list
OS_ARCHS = linux_386 linux_amd64 linux_arm linux_arm64 \
           darwin_386 darwin_amd64 \
           windows_386 windows_amd64
OS_ARCHS_LINUX = $(filter linux_%, $(OS_ARCHS))
OS_ARCHS_MAC = $(filter darwin_%, $(OS_ARCHS))
OS_ARCHS_WIN = $(filter windows_%, $(OS_ARCHS))

APPS = $(shell ls cmd)
RELEASE_BIN = authproxy
CACHE = .cache
mkdir = @mkdir -p $(dir $@)
GO_FILES_NO_TEST := $(shell find . -name "*.go" -not -name "*_test.go")
MD_FILES := $(shell find . -name "*.md")

.PHONY: all
all: dep format build lint test

.PHONY: clean
clean:
	@rm -rf bin dist $(CACHE)
	@find . -type d -name ".snapshots" -exec rm -rf '{}' '+'

.PHONY: dep dep-new dep-update
dep:
	go mod download
	go get -u golang.org/x/tools/cmd/goimports
	go mod tidy -v
	go mod verify
dep-new:
	go get ./...
	go get -u golang.org/x/tools/cmd/goimports
	go mod tidy -v
	go mod verify
dep-update:
	go get ./...
	@go list -m -u -json all | jq -r '. | select(.Indirect != true) | select(.Update != null) | .Path' | while read pkg; do echo "go get -u $$pkg"; go get -u $$pkg; done
	go get -u golang.org/x/tools/cmd/goimports
	go mod tidy -v
	go mod verify

.PHONY: format
format:
	$(MAKE) \
	format/go \
	format/yaml \
	format/md

.PHONY: format/go format/yaml format/md $(MD_FILES)
format/go:
	gofmt -s -w -e -l .
	goimports -w -e -l .
format/yaml:
	@which prettier > /dev/null && prettier --write "**/*.yaml" "**/*.yml" || npx prettier --write "**/*.yaml" "**/*.yml"
format/md:
	$(MAKE) \
	$(MD_FILES)
$(MD_FILES):
	@which markdown-toc > /dev/null && markdown-toc -i $@ || npx markdown-toc -i $@
	@which prettier > /dev/null && prettier --write $@ || npx prettier --write $@

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: test-codecov
test-codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	bash <(curl -s https://codecov.io/bash)

# Compile binaries
.PHONY: build
build:
	$(MAKE) \
	$(sort $(foreach a, $(OS_ARCHS_LINUX), $(addprefix bin/$(VERSION)/$a/, $(APPS)))) \
	$(sort $(foreach a, $(OS_ARCHS_MAC), $(addprefix bin/$(VERSION)/$a/, $(APPS)))) \
	$(sort $(foreach a, $(OS_ARCHS_WIN), $(addsuffix .exe, $(addprefix bin/$(VERSION)/$a/, $(APPS)))))

# Magic for parsing OS, architecture and package name for bin/$(VERSION)/% target
bin_parts = $(subst /, , $@)
bin_os_arch = $(subst _, , $(word 3, $(bin_parts)))
bin_os = $(word 1, $(bin_os_arch))
bin_arch = $(word 2, $(bin_os_arch))
bin_pkg = $(basename $(word 4, $(bin_parts)))

# Compile a binary with the specified OS and architecture. Example:
#   make bin/$(VERSION)/linux_amd64/<app>
.PRECIOUS: bin/$(VERSION)/%
bin/$(VERSION)/%: $(GO_FILES_NO_TEST) go.mod go.sum
	CGO_ENABLED=0 GOOS=$(bin_os) GOARCH=$(bin_arch) go build -trimpath -installsuffix 'static' -ldflags "-X main.version=$(VERSION) -s -w" -o $@ $(PACKAGE)/cmd/$(bin_pkg)

.PRECIOUS: dist/$(VERSION)/$(RELEASE_BIN)_%.tar.gz
dist/$(VERSION)/$(RELEASE_BIN)_%.tar.gz: bin/$(VERSION)/%/$(RELEASE_BIN)
	$(mkdir)
	tar -zcvf $@ --directory="bin/$(VERSION)/$*" $(RELEASE_BIN)

.PRECIOUS: dist/$(VERSION)/$(RELEASE_BIN)_windows_%.tar.gz
dist/$(VERSION)/$(RELEASE_BIN)_windows_%.tar.gz: bin/$(VERSION)/windows_%/$(RELEASE_BIN).exe
	$(mkdir)
	tar -zcvf $@ --directory="bin/$(VERSION)/windows_$*" $(RELEASE_BIN).exe

.PHONY: release
release:
	$(MAKE) \
	$(sort $(foreach a, $(OS_ARCHS_LINUX), github/$(VERSION)/$(RELEASE_BIN)_$a.tar.gz)) \
	$(sort $(foreach a, $(OS_ARCHS_MAC), github/$(VERSION)/$(RELEASE_BIN)_$a.tar.gz)) \
	$(sort $(foreach a, $(OS_ARCHS_WIN), github/$(VERSION)/$(RELEASE_BIN)_$a.tar.gz))

.PHONY: github/$(VERSION)/%
github/$(VERSION)/%: dist/$(VERSION)/% $(CACHE)/$(VERSION)/github-upload-url
	@if [ -z $(GITHUB_API_USERNAME) ]; then echo "GITHUB_API_USERNAME environment variable not set"; exit 1; fi
	@if [ -z $(GITHUB_API_TOKEN) ]; then echo "GITHUB_API_TOKEN environment variable not set"; exit 1; fi
	@curl --request POST \
	--user $$GITHUB_API_USERNAME:$$GITHUB_API_TOKEN \
	--url "$(shell cat $(word 2,$^))?name=$(notdir $<)" \
	--header 'content-type: application/gzip' \
	--data '@$<'

$(CACHE)/$(VERSION)/github-upload-url:
	@if [ -z $(GITHUB_API_USERNAME) ]; then echo "GITHUB_API_USERNAME environment variable not set"; exit 1; fi
	@if [ -z $(GITHUB_API_TOKEN) ]; then echo "GITHUB_API_TOKEN environment variable not set"; exit 1; fi
	$(mkdir)
	@curl --request GET \
    	--url '$(GITHUB_API_URL)/releases/tags/$(TAG)' \
    	--user $(GITHUB_API_USERNAME):$(GITHUB_API_TOKEN) \
    	| jq -r '.upload_url' | sed 's|\(.*/assets\){.*}|\1|' > $@

.PHONY: parallelism
parallelism:
	@echo $(PARALLELISM)

