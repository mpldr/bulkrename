VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go tool dist env | env grep "^GOOS" | env sed 's/[^"]*"\([^"]*\)"/\1/')

build:
ifeq (${GOOS}, windows)
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o br.exe
else
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o br
endif

check:
	gocyclo -over 15 -avg .
	golint -set_exit_status ./...


test:
	go test -v -cover -race ./...

coveralls:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go test -v -covermode=count -coverprofile=coverage.out
	$(shell go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.out -repotoken ${COVERALLS_TOKEN}

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

doc:
	asciidoc -b docbook documentation/manpage.1.txt
	asciidoc documentation/manpage.1.txt
	xsltproc --nonet /etc/asciidoc/docbook-xsl/manpage.xsl documentation/manpage.1.xml

clean:
	$(RM) br.1 documentation/manpage.1.html documentation/manpage.1.xml bulkrename bulkrename.exe br br.exe coverage.out

all: test build doc
