FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.18 as builder

ARG VERSION=v0.0.0

RUN apk add --no-cache build-base libc-dev
RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-batcher with the shared go.mod & go.sum files
COPY ./op-batcher /app/op-batcher
COPY ./op-node /app/op-node
COPY ./op-service /app/op-service
COPY ./op-aws-sdk /app/op-aws-sdk
COPY ./op-plasma /app/op-plasma
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./.git /app/.git

WORKDIR /app/op-batcher

RUN go mod download

ARG TARGETOS TARGETARCH

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
RUN make op-batcher VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.18

# Set user variables
ARG BATCHER_USER_UID=1001
ARG BATCHER_USER_GID=1001
ARG BATCHER_USER=batcher

RUN apk add --no-cache libstdc++ ca-certificates bash \
    && rm -rf /var/cache/apk/* \
    && addgroup -g ${BATCHER_USER_GID} ${BATCHER_USER} \
    && adduser -u ${BATCHER_USER_UID} -G ${BATCHER_USER} -D ${BATCHER_USER} \
    && addgroup ${BATCHER_USER} tty \
    && sed -i -e "s/bin\/sh/bin\/bash/" /etc/passwd

RUN echo "[ ! -z \"\$TERM\" -a -r /etc/motd ] && cat /etc/motd" >> /etc/bash/bashrc

COPY --from=builder /app/op-batcher/bin/op-batcher /usr/local/bin

# Change ownership of the binary to the non-root user
RUN chown ${BATCHER_USER}:${BATCHER_USER} /usr/local/bin/op-batcher

# Set working directory to user's home
WORKDIR /home/${BATCHER_USER}

# Switch to non-root user
USER ${BATCHER_USER}

ENTRYPOINT ["op-batcher"]
