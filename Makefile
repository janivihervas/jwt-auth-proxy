SHELL := /bin/bash
MAKEFLAGS += --no-print-directory --output-sync

BINARY := jwt-auth-proxy
CMD := github.com/janivihervas/$(BINARY)/cmd/$(BINARY)
VERSION ?= $(shell git rev-parse HEAD)

GO_FILES_NO_TEST := $(shell find . -name "*.go" -not -path "./bin/*" -not -path ".git/*" -not -name "*_test.go")
GO_TOOLS := golang.org/x/tools/cmd/goimports \
            github.com/golang/lint/golint \
            github.com/fzipp/gocyclo \
            github.com/kisielk/errcheck \
            github.com/alexkohler/nakedret \
            mvdan.cc/interfacer

.PHONY: all
all: build

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: install
install:
	go get ./...

.PHONY: setup
setup:
	go get -u $(GO_TOOLS)

.PHONY: format
format:
	gofmt -s -w -e -l .
	goimports -w -e -l .

.PHONY: vet golint gocyclo interfacer errcheck nakedret
vet:
	go vet ./...
golint:
	golint -set_exit_status ./...
gocyclo:
	gocyclo -over 12 $(GO_FILES_NO_TEST)
interfacer:
	interfacer ./...
errcheck:
	errcheck -ignoretests ./...
nakedret:
	nakedret ./...
.PHONY: lint
lint:
	@$(MAKE) -j \
	vet \
	golint \
	gocyclo \
	interfacer \
	errcheck \
	nakedret

.PHONY: test
test:
	UPDATE_SNAPSHOTS=true go test -race -cover ./...

.PHONY: test-codecov
test-codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	bash <(curl -s https://codecov.io/bash)

parts = $(subst -, ,$*)
os = $(word 1, $(parts))
arch = $(word 2, $(parts))

.PRECIOUS: bin/$(BINARY).$(VERSION).%/$(BINARY)
bin/$(BINARY).$(VERSION).%/$(BINARY): $(GO_FILES_NO_TEST)
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -installsuffix cgo -o $@ $(CMD)

bin/$(BINARY).$(VERSION).%.tar.gz: bin/$(BINARY).$(VERSION).%/$(BINARY)
	tar -zcvf $@ --directory="bin" $(subst .tar.gz,,$(notdir $@))

.PRECIOUS: bin/$(BINARY).$(VERSION).windows-%/$(BINARY).exe
bin/$(BINARY).$(VERSION).windows-%/$(BINARY).exe: $(GO_FILES_NO_TEST)
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=windows GOARCH=$* go build -installsuffix cgo -o $@ $(CMD)

bin/$(BINARY).$(VERSION).windows-%.tar.gz: bin/$(BINARY).$(VERSION).windows-%/$(BINARY).exe
	tar -zcvf $@ --directory="bin" $(subst .tar.gz,,$(notdir $@))

.PHONY: build
build:
	go install $(CMD)
	@$(MAKE) -j \
	bin/$(BINARY).$(VERSION).linux-amd64/$(BINARY) \
	bin/$(BINARY).$(VERSION).darwin-amd64/$(BINARY) \
	bin/$(BINARY).$(VERSION).windows-amd64/$(BINARY).exe

.PHONY: release
release:
	@$(MAKE) -j \
	bin/$(BINARY).$(VERSION).linux-amd64.tar.gz \
	bin/$(BINARY).$(VERSION).darwin-amd64.tar.gz \
	bin/$(BINARY).$(VERSION).windows-amd64.tar.gz
