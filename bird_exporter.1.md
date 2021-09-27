---
date: 2018-06-20
footer: bird_exporter
header: "bird_exporter's Manual"
layout: page
license: "Licensed under the MIT license"
section: 1
title: BIRD_EXPORTER
---

# NAME

bird_exporter - A protocol state exporter for the BIRD routing daemon to use
with Prometheus

# SYNOPSIS

**bird_exporter** [**OPTIONS**]

# DESCRIPTION

**bird_exporter** is a metric exporter for the BIRD routing daemon to use with
Prometheus. Since **bird_exporter** uses the BIRD Unix socket(s), BIRD needs to
be installed on the same machine as bird_exporter. The user executing
bird_exporter must have read/write permission to access the BIRD Unix sockets.

# OPTIONS

**-bird.ipv4**
    Get protocols from bird (not compatible with **-bird.v2**)

**-bird.ipv6**
    Get protocols from bird6 (not compatible with **-bird.v2**)

**-bird.socket** */path/to/socket*
    Socket to communicate with bird routing daemon

**-bird.socket6** */path/to/socket*
    Socket to communicate with bird6 routing daemon (not compatible with
**-bird.v2**)

**-bird.v2**
    BIRD major version >= 2.0 (multi channel protocols)

**-format.new**
    New metric format (more convenient / generic)

**-proto.bgp**
    Enables metrics for protocol BGP

**-proto.direct**
    Enables metrics for protocol Direct

**-proto.kernel**
    Enables metrics for protocol Kernel

**-proto.ospf**
    Enables metrics for protocol OSPF

**-proto.static**
    Enables metrics for protocol Static

**-proto.babel**
    Enables metrics for protocol Babel

**-version**
    Print version information

**-web.listen-address** *[address]:port*
    Address on which to expose metrics and web interface

**-web.telemetry-path** *path*
    Path under which to expose metrics (default "/metrics")

Version 2.0 of BIRD supports both IPv4 and IPv6 in a single daemon. Since
version 1.1 of **bird_exporter**, it can be used with BIRD 2.0+ using the
**-bird.v2** option. When using this option, **bird_exporter** queries the same
socket for both IPv4 and IPv6. In this mode the IP protocol is determined by
the channel information, and options **-bird.ipv4**, **-bird.ipv6** and
**-bird.socket6** are ignored.

# BIRD CONFIGURATION

To get meaningful uptime information, BIRD needs to be configured to use
ISO-format timestamps:

```
timeformat protocol iso long;
```

# METRIC FORMATS

In version 1.0, a new metric format was introduced. To avoid backwards
incompatibility, the new format is optional and can be enabled by using the
**-format.new** option. The new format handles protocols more generically and
allows for a better query structure. It also adheres more to the Prometheus
metric naming best practices. In both formats protocol specific metrics are
prefixed with the protocol name (e.g. OSPF running metric).

## OLD METRIC FORMAT EXAMPLE

```
bgp4_session_prefix_count_import{name="bgp1"} 600000
bgp6_session_prefix_count_import{name="bgp1"} 50000
ospfv3_running{name="ospf1"} 1
```

## NEW METRIC FORMAT EXAMPLE

```
bird_protocol_prefix_import_count{name="bgp1",proto="BGP",ip_version="4"} 600000
bird_protocol_prefix_import_count{name="bgp1",proto="BGP",ip_version="6"} 50000
bird_ospfv3_running{name="ospf1"} 1
```

# AUTHOR

Daniel Czerwonk <daniel@dan-nrw.de>
