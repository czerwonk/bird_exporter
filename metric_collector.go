package main

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricCollector struct {
	Protocols []ProtocolMetric
}

func NewMetricCollectorForProtocols(p []*protocol.Protocol) *MetricCollector {
	m := make([]ProtocolMetric, 0)
	for _, x := range p {
		m = append(m, NewProtocolMetricFromProtocol(x))
	}

	return &MetricCollector{Protocols: m}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, p := range m.Protocols {
		p.Describe(ch)
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, p := range m.Protocols {
		p.GetMetrics(ch)
	}
}
