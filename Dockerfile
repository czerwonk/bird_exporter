FROM --platform=$BUILDPLATFORM golang AS builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ADD . /go/bird_exporter/
WORKDIR /go/bird_exporter
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	GOARM=${TARGETVARIANT#v} \
	go build -a -installsuffix cgo -o /go/bin/bird_exporter

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash tzdata
WORKDIR /app
COPY --from=builder /go/bin/bird_exporter .
EXPOSE 9324

ENTRYPOINT ["/app/bird_exporter"]
