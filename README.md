# bird_exporter 
[![Build Status](https://travis-ci.org/czerwonk/bird_exporter.svg)](https://travis-ci.org/czerwonk/bird_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/bird_exporter)](https://goreportcard.com/report/github.com/czerwonk/bird_exporter)

Metric exporter for bird routing daemon to use with Prometheus

## Remarks
Since bird_exporter uses the bird unix sockets, bird has to be installed on the same mashine as bird_exporter. Also the user executing bird_exporter must have permission to access the bird socket files. 

### Bird configuration
To get meaningful uptime information bird has to be configured this way:
```
timeformat protocol     iso long;
```

## Important information for users of bird 2.0+
Version 2.0 of bird routing daemon does support IPv4 and IPv6 in one single daemon now. 
For further information see [here](https://gitlab.labs.nic.cz/labs/bird/wikis/transition-notes-to-bird-2). 
Since version 1.1 bird_exporter can be used with bird 2.0+ using the `-bird.v2` parameter. 
When using this parameter bird_exporter queries the same bird socket for IPv4 and IPv6. 
In this mode the IP protocol is determined by the channel information and parameters `-bird.ipv4`, `-bird.ipv6` and `-bird.socket6` are ignored.

## Metric formats
In version 1.0 a new metric format was introduced. 
To prevent a breaking change the new format is optional and can be enabled by using the ```-format.new``` flag.
The new format handles protocols more generic and allows a better query structure.
Also it adheres more to the metric naming best practices.
In both formats protocol specific metrics are prefixed with the protocol name (e.g. OSPF running metric).

This is a short example of the different formats:

### old format
```
bgp4_session_prefix_count_import{name="bgp1"} 600000
bgp6_session_prefix_count_import{name="bgp1"} 50000
ospfv3_running{name="ospf1"} 1
```

### new format
```
bird_protocol_prefix_import_count{name="bgp1",proto="BGP",ip_version="4"} 600000
bird_protocol_prefix_import_count{name="bgp1",proto="BGP",ip_version="6"} 50000
bird_ospfv3_running{name="ospf1"} 1
```

### Default Port
In version 0.7.1 the default port changed to 9324 since port 9200 is the default port of elasticsearch. The new port is now registered in the default port allocation list (https://github.com/prometheus/prometheus/wiki/Default-port-allocations)

### Sockets
In version 0.8 communication to bird changed to sockets. The default socket path is ```/var/run/bird.ctl``` (for bird) and ```/var/run/bird6.ctl``` (for bird6). In case you are using different paths in your installation, the socket path can be specified by usind the ```-bird.socket``` (for bird) and ```-bird.socket6``` (for bird6) flag.

## Install
```
go get -u github.com/czerwonk/bird_exporter
```

## Usage
```
bird_exporter -format.new=true
```

###### BIRD RS Dashboard

https://grafana.com/dashboards/5259

![alt text](https://github.com/openbsod/bird_exporter/blob/master/img/bird_exporter.png)

## Features
* BGP session state
* OSPF neighbor/interface count
* imported / exported / filtered prefix counts / route state changes (BGP, OSPF, Kernel, Static, Device, Direct)
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
