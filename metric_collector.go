package main

import (
	"github.com/czerwonk/bird_exporter/bgp"
	"github.com/czerwonk/bird_exporter/ospf"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricCollector struct {
	collectors []prometheus.Collector
}

func NewMetricCollectorForProtocols(protocols []*protocol.Protocol) *MetricCollector {
	b := make([]*protocol.Protocol, 0)
	o := make([]*protocol.Protocol, 0)

	for _, p := range protocols {
		if p.Proto == protocol.BGP {
			b = append(b, p)
		} else if p.Proto == protocol.OSPF {
			o = append(o, p)
		}
	}

	c := make([]prometheus.Collector, 0)
	if len(b) > 0 {
		c = append(c, bgp.NewCollector(b))
	}
	if len(o) > 0 {
		c = append(c, ospf.NewCollector(o))
	}

	return &MetricCollector{collectors: c}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range m.collectors {
		c.Describe(ch)
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, c := range m.collectors {
		c.Collect(ch)
	}
}
