VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go tool dist banner | head -2 | tail -1 | sed -r 's/[^/]* ([a-z0-9]+)\/[A-Za-z0-9 \/]*/\1/')

build:
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}"

check:
	gocyclo -over 15 -avg .
	golint -set_exit_status ./...


test:
	go test -v -cover

coveralls:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go test -v -covermode=count -coverprofile=coverage.out
	$(shell go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.out -repotoken ${COVERALLS_TOKEN}
