VERSION := $(shell git describe --always --long --dirty)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOEXE := $(shell go env GOEXE)
GOBIN := $(shell go env GOPATH)/bin/
BINARY := br

build:
	go build -ldflags="-s -w -X main.buildVersion=${VERSION}" -trimpath -buildmode=pie -o $(BINARY)_$(GOOS)_$(GOARCH)$(GOEXE)

check:
	gofmt -s -w -l .
	$(GOBIN)golangci-lint run
	$(GOBIN)gocyclo -over 15 -avg -ignore "_test|Godeps|vendor/" .


test:
	go test -v -cover -race ./...

cover:
	go test -v -coverprofile=coverage.out -covermode=atomic -race ./...
	go tool cover -html=coverage.out

show:
	go test -coverprofile=coverage.out -covermode=atomic -race ./...
	go tool cover -html=coverage.out

doc:
	help2man --version-option=--version --help-option=--help --no-discard-stderr ./$(BINARY)_$(GOOS)_$(GOARCH)$(GOEXE) > br.1.gen
	@head -n -11 br.1.gen > br.1
	@echo "The full documentation is available in Wiki-form at GitLab." >> br.1
	@echo "The reason being that I have no idea about TexInfo and I can't be bothered to learn it." >> br.1
	@echo "Until someone writes a texinfo document and commits to maintaining it, the wiki will be your best help." >> br.1
	@echo -e "\nYou may find the aforementioned wiki at https://gitlab.com/poldi1405/bulkrename/-/wikis/home" >> br.1

clean:
	$(RM) -r br* pkg/ releases/

all: test build doc

release: build doc

prepare:
	go get -v github.com/golangci/golangci-lint@v1.33.1 -o $(GOBIN)golangci-lint$(GOEXE)
	go get -v github.com/fzipp/gocyclo/cmd/gocyclo -o $(GOBIN)gocyclo$(GOEXE)
	go mod tidy
