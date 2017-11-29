package bgp

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var exporter map[int]*protocol.GenericProtocolMetricExporter

type BgpMetricExporter struct {
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("bgp4_session")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("bgp6_session")
}

func (*BgpMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
}

func (*BgpMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	exporter[p.IpVersion].Export(p, ch)
}
