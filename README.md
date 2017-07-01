# bird_exporter 
[![Build Status](https://travis-ci.org/czerwonk/bird_exporter.svg)](https://travis-ci.org/czerwonk/bird_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/bird_exporter)](https://goreportcard.com/report/github.com/czerwonk/bird_exporter)

Metric exporter for bird routing daemon to use with Prometheus

## Remarks
Since bird_exporter uses the bird unix sockets, bird has to be installed on the same mashine as bird_exporter. Also the user executing bird_exporter must have permission to access the bird socket files. 

### Bird configuration
To get meaningful uptime information bird has to be configured this way:
```
timeformat protocol "%s";
```

### Default Port
In version 0.7.1 the default port changed to 9324 since port 9200 is the default port of elasticsearch. The new port is now registered in the default port allocation list (https://github.com/prometheus/prometheus/wiki/Default-port-allocations)

### Sockets
In version 0.8 communication to bird changed to sockets. The default socket path is ```/var/run/bird.ctl``` (for bird) and ```/var/run/bird6.ctl``` (for bird6). In case you are using different paths in your installation, the socket path can be specified by usind the ```-bird.socket``` (for bird) and ```-bird.socket6``` (for bird6) flag.

## Install
```
go get -u github.com/czerwonk/bird_exporter
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
