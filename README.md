# bird_exporter 
[![Build Status](https://travis-ci.org/czerwonk/bird_exporter.svg)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/bird_exporter)][goreportcard]

Metric exporter for bird routing daemon to use with Prometheus

## Remarks
this is an early version

Since bird_exporter calls the bird client binary, bird has to be installed on the same mashine as bird_exporter.

To get meaningful uptime information bird has to be configured this way:
```
timeformat protocol "%s";
```


In version 0.7.1 bird_exporter the default port changed to 9324 since port 9200 is the default port of elasticsearch. The new port is now registered in the default port allocation list (https://github.com/prometheus/prometheus/wiki/Default-port-allocations)

## Install
```
go get github.com/czerwonk/bird_exporter
```

## Features
* BGP session state
* imported / exported / filtered prefix counts (BGP, OSPF)
* protocol uptimes (BGP, OSPF)

## Third Party Components
This software uses components of the following projects
* Prometheus Go client library (https://github.com/prometheus/client_golang)

## License
(c) Daniel Czerwonk, 2016. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Bird routing daemon
see http://bird.network.cz/

[travis]: https://travis-ci.org/czerwonk/bird_exporter
[goreportcard]: https://goreportcard.com/report/github.com/czerwonk/bird_exporter
