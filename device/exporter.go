package device

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var exporter map[int]*protocol.GenericProtocolMetricExporter

type DeviceMetricExporter struct {
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("device4")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("device6")
}

func (m *DeviceMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
}

func (m *DeviceMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	exporter[p.IpVersion].Export(p, ch)
}
