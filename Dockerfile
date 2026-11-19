FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine3.24@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

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


FROM gcr.io/distroless/static-debian13:nonroot@sha256:aafe993665bf6200db38a02c11a9351e782a84b6a45e7f33b6695638c88e3931

WORKDIR /app
COPY --from=builder /out/bird_exporter /app/bird_exporter

USER 1000:1000
EXPOSE 9324
ENTRYPOINT ["/app/bird_exporter"]
