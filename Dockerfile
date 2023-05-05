FROM golang as builder
ADD . /go/bird_exporter/
WORKDIR /go/bird_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/bird_exporter

FROM alpine:3.17
RUN apk update && apk upgrade --no-cache && apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=builder /go/bin/bird_exporter .
EXPOSE 9324


ENTRYPOINT ["/app/bird_exporter"]
