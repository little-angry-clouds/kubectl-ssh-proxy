BIN = $(CURDIR)/bin
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN)
	@tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) go get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint
$(BIN)/gopherbadger: PACKAGE=github.com/jpoles1/gopherbadger

GOPHERBADGER = $(BIN)/gopherbadger
all: clean lint test build | $(GOPHERBADGER)
	$(GOPHERBADGER) -md="README.md"

# Build binaries
build:
	go build -a -o bin/kubectl-ssh_proxy cmd/main/main.go
	go build -a -o bin/kube-ssh-proxy-ssh-bin cmd/ssh/*.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean:
	rm bin/* 2> /dev/null

GOLINT = $(BIN)/golint
lint: | $(GOLINT)
	$(GOLINT) ./...

test: fmt vet
	go test -coverprofile cover.out ./...

PLATFORMS := linux-amd64 linux-386 darwin-amd64 darwin-386 windows-amd64 windows-386
temp = $(subst -, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
release: $(PLATFORMS)
$(PLATFORMS):
	@mkdir -p releases; \
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -a -o bin/kubectl-ssh_proxy-$(os)-$(arch) cmd/main/main.go; \
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -a -o bin/kube-ssh-proxy-ssh-bin-$(os)-$(arch) cmd/ssh/*.go; \
	tar -cvzf releases/kubectl-ssh-proxy-$(os)-$(arch).tar.gz bin/*-$(os)-$(arch)
