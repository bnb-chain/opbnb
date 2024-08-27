FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.18 as builder

ARG VERSION=v0.0.0

RUN apk add --no-cache build-base libc-dev
RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-proposer with the shared go.mod & go.sum files
COPY ./op-proposer /app/op-proposer
COPY ./op-node /app/op-node
COPY ./op-service /app/op-service
COPY ./op-plasma /app/op-plasma
COPY ./op-aws-sdk /app/op-aws-sdk
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./.git /app/.git

WORKDIR /app/op-proposer

RUN go mod download

ARG TARGETOS TARGETARCH

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
RUN make op-proposer VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.18

RUN apk add --no-cache libstdc++

COPY --from=builder /app/op-proposer/bin/op-proposer /usr/local/bin

CMD ["op-proposer"]
