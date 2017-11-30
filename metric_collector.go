package main

import (
	"github.com/czerwonk/bird_exporter/ospf"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricExporter interface {
	Describe(ch chan<- *prometheus.Desc)
	Export(p *protocol.Protocol, ch chan<- prometheus.Metric)
}

type MetricCollector struct {
	protocols []*protocol.Protocol
	exporters map[int]MetricExporter
}

func NewMetricCollectorForProtocols(protocols []*protocol.Protocol) *MetricCollector {
	e := map[int]MetricExporter{
		protocol.BGP:    protocol.NewMetricExporter("bgp4_session", "bgp6_session"),
		protocol.Device: protocol.NewMetricExporter("device4", "device6"),
		protocol.Direct: protocol.NewMetricExporter("direct4", "direct6"),
		protocol.Kernel: protocol.NewMetricExporter("kernel4", "kernel6"),
		protocol.OSPF:   &ospf.OspfMetricExporter{},
		protocol.Static: protocol.NewMetricExporter("static4", "static6"),
	}

	return &MetricCollector{protocols: protocols, exporters: e}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range m.exporters {
		v.Describe(ch)
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, p := range m.protocols {
		if p.Proto != protocol.PROTO_UNKNOWN {
			m.exporters[p.Proto].Export(p, ch)
		}
	}
}
