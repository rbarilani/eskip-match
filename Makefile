BUILD_DIR  ?= build
GOX_ARCH   ?= darwin/386 darwin/amd64 linux/amd64 linux/386
GOX_OUTPUT ?= $(BUILD_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}

install:
	go get -t -v ./...
	go get github.com/mitchellh/gox

test:
	go test ./...

test.coverage:
	go test ./... -coverprofile=coverage.txt -covermode=atomic

build:
	gox -output "$(GOX_OUTPUT)" -osarch="$(GOX_ARCH)"

build.clean:
	rm $(BUILD_DIR)/*
