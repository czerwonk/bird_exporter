package main

import (
	"github.com/czerwonk/bird_exporter/bgp"
	"github.com/czerwonk/bird_exporter/device"
	"github.com/czerwonk/bird_exporter/direct"
	"github.com/czerwonk/bird_exporter/kernel"
	"github.com/czerwonk/bird_exporter/ospf"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/czerwonk/bird_exporter/static"
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
		protocol.BGP:    &bgp.BgpMetricExporter{},
		protocol.Device: &device.DeviceMetricExporter{},
		protocol.Direct: &direct.DirectMetricExporter{},
		protocol.Kernel: &kernel.KernelMetricExporter{},
		protocol.OSPF:   &ospf.OspfMetricExporter{},
		protocol.Static: &static.StaticMetricExporter{},
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
