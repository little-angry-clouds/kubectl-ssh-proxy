# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Build binaries
build: fmt vet
	go build -o bin/kubectl-ssh_proxy cmd/main/main.go
	go build -o bin/kube-ssh-proxy-ssh-bin cmd/ssh/main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean:
	rm bin/*

test: build
	golint ./... | grep -v "should not use dot imports"
	go test -coverprofile cover.out \
		github.com/little-angry-clouds/kubectl-ssh-proxy/cmd/main
	@gopherbadger -md="README.md" > /dev/null
