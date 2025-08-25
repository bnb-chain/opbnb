FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.18 as builder

ARG VERSION=v0.0.0

RUN apk add --no-cache build-base libc-dev wget
RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-node with the shared go.mod & go.sum files
COPY ./op-node /app/op-node
COPY ./op-chain-ops /app/op-chain-ops
COPY ./op-service /app/op-service
COPY ./op-plasma /app/op-plasma
COPY ./op-conductor /app/op-conductor
COPY ./op-aws-sdk /app/op-aws-sdk
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./.git /app/.git

WORKDIR /app/op-node
RUN go mod download

ARG TARGETOS
ARG TARGETARCH

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
RUN make op-node VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.18

# Set user variables
ARG NODE_USER_UID=1001
ARG NODE_USER_GID=1001
ARG NODE_USER=op-node

RUN apk add --no-cache libstdc++ ca-certificates bash curl \
    && rm -rf /var/cache/apk/* \
    && addgroup -g ${NODE_USER_GID} ${NODE_USER} \
    && adduser -u ${NODE_USER_UID} -G ${NODE_USER} -D ${NODE_USER} \
    && addgroup ${NODE_USER} tty \
    && sed -i -e "s/bin\/sh/bin\/bash/" /etc/passwd

RUN echo "[ ! -z \"\$TERM\" -a -r /etc/motd ] && cat /etc/motd" >> /etc/bash/bashrc

COPY --from=builder /app/op-node/bin/op-node /usr/local/bin

# Change ownership of the binary to the non-root user
RUN chown ${NODE_USER}:${NODE_USER} /usr/local/bin/op-node

# Set working directory to user's home
WORKDIR /home/${NODE_USER}

# Switch to non-root user
USER ${NODE_USER}

CMD ["op-node"]
