package main

import (
	"github.com/czerwonk/bird_exporter/metrics"
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
	l := &metrics.LegacyLabelStrategy{}
	e := map[int]MetricExporter{
		protocol.BGP:    metrics.NewMetricExporter("bgp4_session", "bgp6_session", l),
		protocol.Device: metrics.NewMetricExporter("device4", "device6", l),
		protocol.Direct: metrics.NewMetricExporter("direct4", "direct6", l),
		protocol.Kernel: metrics.NewMetricExporter("kernel4", "kernel6", l),
		protocol.OSPF:   ospf.NewExporter(l),
		protocol.Static: metrics.NewMetricExporter("static4", "static6", l),
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
