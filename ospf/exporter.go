package ospf

import (
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var descriptions map[int]*desc

type desc struct {
	runningDesc *prometheus.Desc
}

type OspfMetricExporter struct {
	genericExporter *metrics.ProtocolMetricExporter
}

func init() {
	descriptions = make(map[int]*desc)
	descriptions[4] = getDesc("ospf")
	descriptions[6] = getDesc("ospfv3")
}

func NewExporter(labelStrategy metrics.LabelStrategy) *OspfMetricExporter {
	return &OspfMetricExporter{genericExporter: metrics.NewMetricExporter("ospf", "ospfv3", labelStrategy)}
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	return d
}

func (m *OspfMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	m.genericExporter.Describe(ch)
	ch <- descriptions[4].runningDesc
	ch <- descriptions[6].runningDesc
}

func (m *OspfMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	m.genericExporter.Export(p, ch)
	ch <- prometheus.MustNewConstMetric(descriptions[p.IpVersion].runningDesc, prometheus.GaugeValue, p.Attributes["running"], p.Name)
}
