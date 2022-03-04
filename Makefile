CWD     := $(shell pwd)
NAME    := tpl
VERSION := 1.0.0
ARCH    := $(shell uname -m)

ifeq (${ARCH},x86_64)
ARCH    := amd64
endif

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitBranch=$(shell git describe --all)' \
           -X 'main.BuildGitRev=$(shell git rev-list --count HEAD)' \
           -X 'main.BuildGitCommit=$(shell git rev-parse HEAD)' \
           -X 'main.BuildDate=$(shell date -u -R)'

export GO111MODULE=on

default:
	@ echo "no default target for Makefile"

clean:
	rm -rf $(NAME) ./_releases ./_build

fmt:
	go fmt ./...

lint:
	go vet ./...

build-all: \
    build-linux \
    build-darwin \
    build-windows

build: build-$(shell go env GOOS)

build-linux: clean fmt
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) \
		go build -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-linux-$(ARCH)

build-darwin: clean fmt
	GOOS=darwin GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-darwin-amd64

build-windows: clean fmt
	GOOS=windows GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-windows-amd64.exe

sha256sum: build-all
	@ for f in $(shell ls ./_releases); do \
		cd $(CWD)/_releases; sha256sum "$$f" >> $$f.sha256; \
	done

release: build-all sha256sum
