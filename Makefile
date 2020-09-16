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
	go test -v -cover

coveralls:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go test -v -covermode=count -coverprofile=coverage.out
	$(shell go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.out -repotoken ${COVERALLS_TOKEN}

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
