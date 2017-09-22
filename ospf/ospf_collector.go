package ospf

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var descriptions map[int]*desc
var exporter map[int]*protocol.GenericProtocolMetricExporter

type desc struct {
	runningDesc *prometheus.Desc
}

type OspfCollector struct {
	protocols []*protocol.Protocol
}

func init() {
	exporter = make(map[int]*protocol.GenericProtocolMetricExporter)
	exporter[4] = protocol.NewGenericProtocolMetricExporter("ospf")
	exporter[6] = protocol.NewGenericProtocolMetricExporter("ospfv3")

	descriptions = make(map[int]*desc)
	descriptions[4] = getDesc("ospf")
	descriptions[6] = getDesc("ospfv3")
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	return d
}

func NewCollector(p []*protocol.Protocol) prometheus.Collector {
	return &OspfCollector{protocols: p}
}

func (m *OspfCollector) Describe(ch chan<- *prometheus.Desc) {
	exporter[4].Describe(ch)
	exporter[6].Describe(ch)
	ch <- descriptions[4].runningDesc
	ch <- descriptions[6].runningDesc
}

func (m *OspfCollector) Collect(ch chan<- prometheus.Metric) {
	for _, p := range m.protocols {
		exporter[p.IpVersion].Export(p, ch)
		ch <- prometheus.MustNewConstMetric(descriptions[p.IpVersion].runningDesc, prometheus.GaugeValue, p.Attributes["running"], p.Name)
	}
}
