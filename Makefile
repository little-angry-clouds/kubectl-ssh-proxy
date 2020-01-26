BIN = $(CURDIR)/bin
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN)
	@tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) go get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

# Build binaries
build: fmt vet
	go build -o bin/kubectl-ssh_proxy cmd/main/*.go
	go build -o bin/kube-ssh-proxy-ssh-bin cmd/ssh/*.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean:
	rm bin/*

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOLINT = $(BIN)/golint
lint: | $(GOLINT)
	$(GOLINT) ./...

test: fmt vet lint
	go test -coverprofile cover.out \
		github.com/little-angry-clouds/kubectl-ssh-proxy/cmd/main
	gopherbadger -md="README.md"
