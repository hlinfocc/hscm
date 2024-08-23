export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: hscm


fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

vet:
	go vet ./...


hscm:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/hscm ./

test: gotest

gotest:
	go test -v --cover ./...

	
clean:
	rm -f ./bin/hscm
