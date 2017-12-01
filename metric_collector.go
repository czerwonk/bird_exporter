package main

import (
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/czerwonk/bird_exporter/ospf"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricCollector struct {
	protocols        []*protocol.Protocol
	exporters        map[int][]metrics.MetricExporter
	enabledProtocols int
}

func NewMetricCollectorForProtocols(protocols []*protocol.Protocol, newFormat bool, enabledProtocols int) *MetricCollector {
	var e map[int][]metrics.MetricExporter

	if newFormat {
		e = exportersForDefault()
	} else {
		e = exportersForLegacy()
	}

	return &MetricCollector{protocols: protocols, exporters: e, enabledProtocols: enabledProtocols}
}

func exportersForLegacy() map[int][]metrics.MetricExporter {
	l := &metrics.LegacyLabelStrategy{}

	return map[int][]metrics.MetricExporter{
		protocol.BGP:    []metrics.MetricExporter{metrics.NewMetricExporter("bgp4_session", "bgp6_session", l)},
		protocol.Device: []metrics.MetricExporter{metrics.NewMetricExporter("device4", "device6", l)},
		protocol.Direct: []metrics.MetricExporter{metrics.NewMetricExporter("direct4", "direct6", l)},
		protocol.Kernel: []metrics.MetricExporter{metrics.NewMetricExporter("kernel4", "kernel6", l)},
		protocol.OSPF:   []metrics.MetricExporter{metrics.NewMetricExporter("ospf", "ospfv3", l), ospf.NewExporter("")},
		protocol.Static: []metrics.MetricExporter{metrics.NewMetricExporter("static4", "static6", l)},
	}
}

func exportersForDefault() map[int][]metrics.MetricExporter {
	l := &metrics.DefaultLabelStrategy{}
	e := metrics.NewGenericProtocolMetricExporter("bird_protocol", true, l)

	return map[int][]metrics.MetricExporter{
		protocol.BGP:    []metrics.MetricExporter{e},
		protocol.Device: []metrics.MetricExporter{e},
		protocol.Direct: []metrics.MetricExporter{e},
		protocol.Kernel: []metrics.MetricExporter{e},
		protocol.OSPF:   []metrics.MetricExporter{e, ospf.NewExporter("bird_")},
		protocol.Static: []metrics.MetricExporter{e},
	}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range m.exporters {
		for _, e := range v {
			e.Describe(ch)
		}
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, p := range m.protocols {
		if p.Proto == protocol.PROTO_UNKNOWN || (m.enabledProtocols & p.Proto != p.Proto) {
			continue
		}

		for _, e := range m.exporters[p.Proto] {
			e.Export(p, ch)
		}
	}
}
