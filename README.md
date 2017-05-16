# bird_exporter [![Build Status](https://travis-ci.org/czerwonk/bird_exporter.svg)][travis]
Metric exporter for bird routing daemon to use with Prometheus

# Remarks
this is an early version

Since bird_exporter calls the bird client binary, bird has to be installed on the same mashine as bird_exporter.

To get meaningful uptime information bird has to be configured this way:
```
timeformat protocol "%s";
```

# Install
```
go get github.com/czerwonk/bird_exporter
```

# Features
* BGP session state
* imported / exported prefix counts (BGP, OSPF)
* protocol uptimes (BGP, OSPF)

# Prometheus
see https://prometheus.io/

# Bird routing daemon
see http://bird.network.cz/

[travis]: https://travis-ci.org/czerwonk/bird_exporter
