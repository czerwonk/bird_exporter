package protocol

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProtocolMetricExporter struct {
	ipv4Exporter *GenericProtocolMetricExporter
	ipv6Exporter *GenericProtocolMetricExporter
}

func NewMetricExporter(prefixIpv4, prefixIpv6 string) *ProtocolMetricExporter {
	return &ProtocolMetricExporter{
		ipv4Exporter: NewGenericProtocolMetricExporter(prefixIpv4),
		ipv6Exporter: NewGenericProtocolMetricExporter(prefixIpv6),
	}
}

func (e *ProtocolMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	e.ipv4Exporter.Describe(ch)
	e.ipv6Exporter.Describe(ch)
}

func (e *ProtocolMetricExporter) Export(p *Protocol, ch chan<- prometheus.Metric) {
	if p.IpVersion == 4 {
		e.ipv4Exporter.Export(p, ch)
	} else {
		e.ipv6Exporter.Export(p, ch)
	}
}
