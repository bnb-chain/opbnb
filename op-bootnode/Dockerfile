FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.18 as builder

ARG VERSION=v0.0.0

RUN apk add --no-cache build-base libc-dev
RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-bootnode with the shared go.mod & go.sum files
COPY ./op-bootnode /app/op-bootnode
COPY ./op-node /app/op-node
COPY ./op-service /app/op-service
COPY ./op-plasma /app/op-plasma
COPY ./op-conductor /app/op-conductor
COPY ./op-aws-sdk /app/op-aws-sdk
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./.git /app/.git

WORKDIR /app/op-bootnode

RUN go mod download

ARG TARGETOS TARGETARCH

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
RUN make op-bootnode VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.18

# Set user variables
ARG BOOTNODE_USER_UID=1001
ARG BOOTNODE_USER_GID=1001
ARG BOOTNODE_USER=bootnode

RUN apk add --no-cache libstdc++ ca-certificates bash \
    && rm -rf /var/cache/apk/* \
    && addgroup -g ${BOOTNODE_USER_GID} ${BOOTNODE_USER} \
    && adduser -u ${BOOTNODE_USER_UID} -G ${BOOTNODE_USER} -D ${BOOTNODE_USER} \
    && addgroup ${BOOTNODE_USER} tty \
    && sed -i -e "s/bin\/sh/bin\/bash/" /etc/passwd

RUN echo "[ ! -z \"\$TERM\" -a -r /etc/motd ] && cat /etc/motd" >> /etc/bash/bashrc

COPY --from=builder /app/op-bootnode/bin/op-bootnode /usr/local/bin

# Change ownership of the binary to the non-root user
RUN chown ${BOOTNODE_USER}:${BOOTNODE_USER} /usr/local/bin/op-bootnode

# Set working directory to user's home
WORKDIR /home/${BOOTNODE_USER}

# Switch to non-root user
USER ${BOOTNODE_USER}

ENTRYPOINT ["op-bootnode"]
