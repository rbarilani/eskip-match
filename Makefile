BUILD_DIR  ?= build
GOX_ARCH   ?= darwin/386 darwin/amd64 linux/amd64 linux/386
GOX_OUTPUT ?= $(BUILD_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}
GOMODULE   ?= GO111MODULE=on

default: install test build

install:
	$(GOMODULE) go get github.com/mitchellh/gox

test:
	$(GOMODULE) go test ./...

test.coverage:
	$(GOMODULE) go test ./... -coverprofile=coverage.txt -covermode=atomic

build: build.clean
	$(GOMODULE) gox -output "$(GOX_OUTPUT)" -osarch="$(GOX_ARCH)"

build.clean:
	if [ -d $(BUILD_DIR) ]; then rm $(BUILD_DIR)/* ; fi
