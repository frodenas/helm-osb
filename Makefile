GO     := GO15VENDOREXPERIMENT=1 go
GINKGO := ginkgo
pkgs   = $(shell $(GO) list ./... | grep -v /vendor/)

all: format build test

deps:
	@$(GO) get github.com/onsi/ginkgo/ginkgo
	@$(GO) get github.com/onsi/gomega

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

test: deps
	@echo ">> running tests"
	@$(GINKGO) -r -race .

.PHONY: all deps format style vet test
