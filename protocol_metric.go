package main

import (
	"github.com/czerwonk/bird_exporter/bgp"
	"github.com/czerwonk/bird_exporter/ospf"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type ProtocolMetric interface {
	Describe(ch chan<- *prometheus.Desc)
	GetMetrics(ch chan<- prometheus.Metric)
}

func NewProtocolMetricFromProtocol(p *protocol.Protocol) ProtocolMetric {
	if p.Proto == protocol.BGP {
		return &bgp.BgpMetric{Protocol: p}
	}

	if p.Proto == protocol.OSPF {
		return &ospf.OspfMetric{Protocol: p}
	}

	return nil
}
