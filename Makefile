NAME    := tpl
LDFLAGS := -s -w

export GO111MODULE=on

default:
	@ echo "no default target for Makefile"

clean:
	rm -rf $(NAME) ./_releases ./_build

fmt:
	go fmt ./...

lint:
	go vet ./...

deps-update:
	go get -u ./...

build-all: \
    build-linux \
    build-darwin

build: build-$(shell go env GOOS)

build-linux: clean fmt
	GOOS=linux GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-linux-arm64
	GOOS=linux GOARCH=riscv64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-linux-riscv64

build-darwin: clean fmt
	GOOS=darwin GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-darwin-arm64

release: build-all
