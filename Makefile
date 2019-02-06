SHELL := /bin/bash
MAKEFLAGS += --no-print-directory --output-sync

BINARY := oidc-go
CMD := github.com/janivihervas/$(BINARY)/cmd/$(BINARY)
VERSION ?= $(shell git rev-parse HEAD)

GO_FILES_NO_TEST := `find . -name "*.go" -not -name "*_test.go"`
GO_TOOLS := golang.org/x/lint/golint \
			golang.org/x/tools/cmd/goimports \
			\
			github.com/alexkohler/nakedret \
			github.com/fzipp/gocyclo \
			github.com/kisielk/errcheck \
			github.com/mdempsky/unconvert \
			\
			gitlab.com/opennota/check/cmd/structcheck \
			gitlab.com/opennota/check/cmd/varcheck \
			\
			honnef.co/go/tools/cmd/staticcheck \

.PHONY: all
all: format build lint test

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: install install-new install-update
install:
	go mod download
	go get -u $(GO_TOOLS)
	go mod tidy -v
	go mod verify
install-new:
	go get ./...
	go get -u $(GO_TOOLS)
	go mod tidy -v
	go mod verify
install-update:
	go get -u ./...
	go get -u $(GO_TOOLS)
	go mod tidy -v
	go mod verify


format:
	gofmt -s -w -e -l .
	goimports -w -e -l .

.PHONY: vet golint
vet:
	go vet ./...
golint:
	golint -set_exit_status ./...

.PHONY: nakedret gocyclo errcheck unconvert
nakedret:
	nakedret ./...
gocyclo:
	gocyclo -over 14 $(GO_FILES_NO_TEST)
errcheck:
	errcheck -ignoretests ./...
unconvert:
	unconvert ./...

.PHONY: structcheck varcheck
structcheck:
	structcheck ./...
varcheck:
	varcheck ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint:
	@$(MAKE) -j \
	vet \
	golint \
	\
	nakedret \
	gocyclo \
	errcheck \
	unconvert \
	\
	structcheck \
	varcheck \
	\
	staticcheck


.PHONY: test
test:
	go test -race -cover ./...

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
