package ospf

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/czerwonk/bird_exporter/metrics"
)

type desc struct {
	runningDesc *prometheus.Desc
}

type ospfMetricExporter struct {
	descriptions map[int]*desc
}

func NewExporter(prefix string) metrics.MetricExporter {
	d := make(map[int]*desc)
	d[4] = getDesc(prefix+"ospf")
	d[6] = getDesc(prefix+"ospfv3")

	return &ospfMetricExporter{descriptions: d}
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	return d
}

func (m *ospfMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.descriptions[4].runningDesc
	ch <- m.descriptions[6].runningDesc
}

func (m *ospfMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.descriptions[p.IpVersion].runningDesc, prometheus.GaugeValue, p.Attributes["running"], p.Name)
}
