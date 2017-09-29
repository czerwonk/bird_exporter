package bgp

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var exporter map[int]*protocol.GenericProtocolMetricExporter

type BgpCollector struct {
	protocols []*protocol.Protocol
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("bgp4_session")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("bgp6_session")
}

func NewCollector(p []*protocol.Protocol) prometheus.Collector {
	return &BgpCollector{protocols: p}
}

func (m *BgpCollector) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
}

func (m *BgpCollector) Collect(ch chan<- prometheus.Metric) {
	for _, p := range m.protocols {
		exporter[p.IpVersion].Export(p, ch)
	}
}
