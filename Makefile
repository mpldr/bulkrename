VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOEXE := $(shell go env GOEXE)
GOBIN := $(shell go env GOPATH)/bin/
BINARY := br

build:
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o $(BINARY)_$(GOOS)_$(GOARCH)$(GOEXE)

check:
	gofumpt -s -w -l .
	$(GOBIN)golangci-lint run

test:
	go test -v -cover -race ./...

cover:
	go test -v -coverprofile=coverage.out -covermode=atomic -race ./...
	go tool cover -html=coverage.out

doc:
	asciidoc -b docbook documentation/br.1.txt
	asciidoc documentation/br.1.txt
	xsltproc --nonet /etc/asciidoc/docbook-xsl/manpage.xsl documentation/br.1.xml
	mv br.1 documentation/

clean:
	$(RM) -r br* pkg/ releases/ bulkrename* docs

all: test build doc

ci:
	make GOOS=linux build
	make GOOS=windows build
	make GOOS=darwin build

release: clean doc
	mkdir docs
	mv documentation/br.1 documentation/br.1.html docs/
	env GOOS=linux make build
	mv br_linux_amd64 br
	tar c docs br | zstd -19 - -o bulkrename.tar.zst
	env GOOS=windows make build
	mv br_windows_amd64.exe br.exe
	zip -r bulkrename_windows.zip docs br.exe
	env GOOS=darwin make build
	mv br_darwin_amd64 br
	tar c docs br | xz -9 > bulkrename_darwin.tar.xz

prepare:
	@go get -v github.com/golangci/golangci-lint@v1.33.1
	@go get -v github.com/fzipp/gocyclo/cmd/gocyclo
	@go get -v golang.org/x/lint/golint
	@go mod tidy
