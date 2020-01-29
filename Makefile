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

# Build binaries
build: clean test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/kubectl-ssh_proxy cmd/main/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/kube-ssh-proxy-ssh-bin cmd/ssh/*.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean:
	rm bin/*

GOLINT = $(BIN)/golint
lint: | $(GOLINT)
	$(GOLINT) -set_exit_status ./...

GOPHERBADGER = $(BIN)/gopherbadger
test: fmt vet | $(GOPHERBADGER)
	go test -coverprofile cover.out \
		github.com/little-angry-clouds/kubectl-ssh-proxy/cmd/main
	$(GOPHERBADGER) -md="README.md"
