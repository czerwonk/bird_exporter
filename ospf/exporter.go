package ospf

import (
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/czerwonk/bird_exporter/client"
)

type desc struct {
	runningDesc *prometheus.Desc
	interfaceCountDesc *prometheus.Desc
	neighborCountDesc *prometheus.Desc
	neighborAdjacentCountDesc *prometheus.Desc
}

type ospfMetricExporter struct {
	descriptions map[string]*desc
	client *client.BirdClient
}

func NewExporter(prefix string, client *client.BirdClient) metrics.MetricExporter {
	d := make(map[string]*desc)
	d["4"] = getDesc(prefix + "ospf")
	d["6"] = getDesc(prefix + "ospfv3")

	return &ospfMetricExporter{descriptions: d, client: client}
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	labels = append(labels, "area")
	d.interfaceCountDesc = prometheus.NewDesc(prefix+"_interface_count", "Number of interfaces in the area", labels, nil)
	d.neighborCountDesc = prometheus.NewDesc(prefix+"_neighbor_count", "Number of neighbors in the area", labels, nil)
	d.neighborAdjacentCountDesc = prometheus.NewDesc(prefix+"_neighbor_adjacent_count", "Number of adjacent neighbors in the area", labels, nil)

	return d
}

func (m *ospfMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.descriptions["4"].runningDesc
	ch <- m.descriptions["6"].runningDesc
}

func (m *ospfMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.descriptions[p.IpVersion].runningDesc, prometheus.GaugeValue, p.Attributes["running"], p.Name)
}
