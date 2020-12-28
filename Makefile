VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go tool dist env | env grep "^GOOS" | env sed 's/[^"]*"\([^"]*\)"/\1/')

build:
ifeq (${GOOS}, windows)
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o br.exe
else
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o br
endif

check:
	gofmt -s -w -l .
	golangci-lint run
	gocyclo -over 15 -avg .


test:
	go test -v -cover -race ./...

cover:
	go test -v -coverprofile=coverage.out -covermode=atomic -race ./...

show:
	go test -coverprofile=coverage.out -covermode=atomic -race ./...
	go tool cover -html=coverage.out

doc:
	asciidoc -b docbook documentation/manpage.1.txt
	asciidoc documentation/manpage.1.txt
	xsltproc --nonet /etc/asciidoc/docbook-xsl/manpage.xsl documentation/manpage.1.xml

clean:
	$(RM) -r br.1 br.1.gz documentation/manpage.1.html documentation/manpage.1.xml bulkrename bulkrename.exe br br.exe pkg/ releases/

all: test build doc

release: clean doc
	mkdir -p pkg/usr/local/bin/ pkg/usr/share/man/man1/ pkg/windows/ releases/
	gzip -cf br.1 > pkg/usr/share/man/man1/br.1.gz
	env GOOS=linux go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o pkg/usr/local/bin/br
	sed -i 's/VERSION/${VERSION}/' .nfpm.yml
	nfpm pkg --packager deb --target ./releases/ --config .nfpm.yml
	nfpm pkg --packager rpm --target ./releases/ --config .nfpm.yml
	nfpm pkg --packager apk --target ./releases/ --config .nfpm.yml
	sed -i 's/${VERSION}/VERSION/' .nfpm.yml
	cd pkg/; tar -c usr/ | xz > ../releases/bulkrename_linux_${VERSION}.tar.xz; cd ..
	env GOOS=windows go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o pkg/windows/br.exe
	cp documentation/manpage.1.html pkg/windows/user_guide.html
	cd pkg/windows/; zip ../../releases/windows_x64_${VERSION}.zip *; cd ../../

prepare:
	go get -v github.com/golangci/golangci-lint@v1.33.1
	go get -v github.com/fzipp/gocyclo/cmd/gocyclo
	go mod tidy
