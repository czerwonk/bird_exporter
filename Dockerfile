from golang:1.10 as builder
arg CMD
run wget -o/dev/null -O/usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
        chmod +x /usr/local/bin/dep
workdir ${GOPATH}/src/github.com/czerwonk/bird_exporter
copy . .
run make deps build && cp bird_exporter /bird_exporter

from debian:stretch-slim
copy --from=builder /bird_exporter /bird_exporter
entrypoint ["/bird_exporter"]
