# bird_exporter [![Build Status](https://travis-ci.org/czerwonk/bird_exporter.svg)][travis]
Metric exporter for bird routing daemon to use with Prometheus

# Remarks
this is an early version

# Install
```
go get github.com/czerwonk/bird_exporter
```

# Features
* BGP session state
* imported / exported prefix counts (BGP, OSPF)
* protocol uptimes (BGP, OSPF)

# Future plans
* systemd unit

# Prometheus
see https://prometheus.io/

# Bird routing daemon
see http://bird.network.cz/

[travis]: https://travis-ci.org/czerwonk/bird_exporter
