VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go tool dist banner | head -2 | tail -1 | sed -r 's/[^/]* ([a-z0-9]+)\/[A-Za-z0-9 \/]*/\1/')

build:
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}"
