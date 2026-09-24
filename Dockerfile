FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine3.24@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS builder
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

FROM gcr.io/distroless/static-debian13:nonroot@sha256:2293b36c7c9082bf4115aab724b4d2cddec82c8eba39bf27ac0517e159acf150
WORKDIR /app
COPY --from=builder /out/bird_exporter /app/bird_exporter
USER 1000:1000
EXPOSE 9324
ENTRYPOINT ["/app/bird_exporter"]
