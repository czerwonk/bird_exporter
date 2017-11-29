package kernel

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var exporter map[int]*protocol.GenericProtocolMetricExporter

type KernelMetricExporter struct {
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("kernel4")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("kernel6")
}

func (*KernelMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
}

func (*KernelMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	exporter[p.IpVersion].Export(p, ch)
}
