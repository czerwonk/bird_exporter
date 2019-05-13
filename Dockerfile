FROM golang as builder
RUN go get -d -v github.com/czerwonk/bird_exporter
WORKDIR /go/src/github.com/czerwonk/bird_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=builder /go/src/github.com/czerwonk/bird_exporter/app bird_exporter
EXPOSE 9324

ENTRYPOINT ["/app/bird_exporter"]
