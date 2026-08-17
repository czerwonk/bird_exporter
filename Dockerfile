FROM --platform=$BUILDPLATFORM golang:1.26.6-alpine3.24@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION=dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY client ./client
COPY metrics ./metrics
COPY parser ./parser
COPY protocol ./protocol

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	GOARM=${TARGETVARIANT#v} \
	go build -trimpath \
	-ldflags="-s -w -X main.version=${VERSION}" \
	-o /out/bird_exporter .

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

RUN apk --no-cache add ca-certificates tzdata \
	&& addgroup -S -g 1000 bird-exporter \
	&& adduser -S -D -H -u 1000 -G bird-exporter bird-exporter

WORKDIR /app
COPY --from=builder /out/bird_exporter /app/bird_exporter

USER 1000:1000
EXPOSE 9324
ENTRYPOINT ["/app/bird_exporter"]
