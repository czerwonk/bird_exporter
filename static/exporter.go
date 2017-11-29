package static

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var exporter map[int]*protocol.GenericProtocolMetricExporter

type StaticMetricExporter struct {
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("static4")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("static6")
}

func (m *StaticMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
}

func (m *StaticMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	exporter[p.IpVersion].Export(p, ch)
}
