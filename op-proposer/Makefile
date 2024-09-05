GITCOMMIT ?= $(shell git rev-parse HEAD)
GITDATE ?= $(shell git show -s --format='%ct')
# Find the github tag that points to this commit. If none are found, set the version string to "untagged"
# Prioritizes release tag, if one exists, over tags suffixed with "-rc"
VERSION ?= $(shell tags=$$(git tag --points-at $(GITCOMMIT) | grep '^op-proposer/' | sed 's/op-proposer\///' | sort -V); \
             preferred_tag=$$(echo "$$tags" | grep -v -- '-rc' | tail -n 1); \
             if [ -z "$$preferred_tag" ]; then \
                 if [ -z "$$tags" ]; then \
                     echo "untagged"; \
                 else \
                     echo "$$tags" | tail -n 1; \
                 fi \
             else \
                 echo $$preferred_tag; \
             fi)

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.Version=$(VERSION)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

op-proposer:
ifeq ($(TARGETARCH),arm64)
	wget https://musl.cc/aarch64-linux-musl-cross.tgz
	tar zxf aarch64-linux-musl-cross.tgz
	export PATH=$$PATH:/app/op-proposer/aarch64-linux-musl-cross/bin/ && env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) CC=aarch64-linux-musl-gcc CGO_ENABLED=1 go build -v $(LDFLAGS) -o ./bin/op-proposer ./cmd
else
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/op-proposer ./cmd
endif

clean:
	rm bin/op-proposer

test:
	go test -v ./...

.PHONY: \
	clean \
	op-proposer \
	test
