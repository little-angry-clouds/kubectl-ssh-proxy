# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Build manager binary
manager: fmt vet
	go build -o bin/kubectl-ssh-proxy cmd/main/main.go
	go build -o bin/kubectl-ssh-proxy-ssh-bin cmd/ssh/main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean:
	rm bin/*
