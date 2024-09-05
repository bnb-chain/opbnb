GITCOMMIT ?= $(shell git rev-parse HEAD)
GITDATE ?= $(shell git show -s --format='%ct')
VERSION ?= v0.0.0

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.Version=$(VERSION)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

op-bootnode:
ifeq ($(TARGETARCH),arm64)
	wget https://musl.cc/aarch64-linux-musl-cross.tgz
	tar zxf aarch64-linux-musl-cross.tgz
	export PATH=$$PATH:/app/op-bootnode/aarch64-linux-musl-cross/bin/ && env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) CC=aarch64-linux-musl-gcc CGO_ENABLED=1 go build -v $(LDFLAGS) -o ./bin/op-bootnode ./cmd
else
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/op-bootnode ./cmd
endif

clean:
	rm -f bin/op-bootnode

test:
	go test -v ./...

.PHONY: \
	op-bootnode \
	clean \
	test
