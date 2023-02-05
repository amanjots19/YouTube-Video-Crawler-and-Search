APP_NAME=yt-fetch
GO_BUILD_DIR=build
VERSION?=$(shell (git describe --tags --exact-match 2> /dev/null || git rev-parse HEAD) | sed "s/^v//")
all: build

build:
	mkdir -p $(GO_BUILD_DIR)
	go build -v -ldflags="-s -w -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...

run:
	$(GO_BUILD_DIR)/$(APP_NAME)

clean:
	rm -r $(GO_BUILD_DIR)