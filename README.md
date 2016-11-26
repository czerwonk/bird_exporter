# go-bird_bgp_exporter
BGP session state exporter for bird routing daemon to use with Prometheus

# Remarks
this is an early version

# Future plans
* Port/IP binding should be customizable
* Atomatic detection of birdcl if birdc is not available (CentOS EPEL)
* Prefix count per session (received, advertised)
* systemd unit
* Support for go install 

# Prometheus
see https://prometheus.io/

# Bird routing daemon
see http://bird.network.cz/
