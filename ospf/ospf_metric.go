package ospf

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

var descriptions map[int]*desc

type desc struct {
	upDesc          *prometheus.Desc
	importCountDesc *prometheus.Desc
	exportCountDesc *prometheus.Desc
	filterCountDesc *prometheus.Desc
	uptimeDesc      *prometheus.Desc
	runningDesc     *prometheus.Desc
}

type OspfMetric struct {
	Protocol *protocol.Protocol
}

func init() {
	descriptions = make(map[int]*desc)
	descriptions[4] = getDesc("ospf")
	descriptions[6] = getDesc("ospfv3")
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.upDesc = prometheus.NewDesc(prefix+"_up", "Protocol is up", labels, nil)
	d.importCountDesc = prometheus.NewDesc(prefix+"_prefix_count_import", "Number of imported routes", labels, nil)
	d.exportCountDesc = prometheus.NewDesc(prefix+"_prefix_count_export", "Number of exported routes", labels, nil)
	d.filterCountDesc = prometheus.NewDesc(prefix+"_prefix_count_filter", "Number of filtered routes", labels, nil)
	d.uptimeDesc = prometheus.NewDesc(prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	return d
}

func (m *OspfMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- descriptions[m.Protocol.IpVersion].upDesc
	ch <- descriptions[m.Protocol.IpVersion].importCountDesc
	ch <- descriptions[m.Protocol.IpVersion].exportCountDesc
	ch <- descriptions[m.Protocol.IpVersion].filterCountDesc
	ch <- descriptions[m.Protocol.IpVersion].uptimeDesc
	ch <- descriptions[m.Protocol.IpVersion].runningDesc
}

func (m *OspfMetric) GetMetrics(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].upDesc, prometheus.GaugeValue, float64(m.Protocol.Up), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].importCountDesc, prometheus.GaugeValue, float64(m.Protocol.Imported), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].exportCountDesc, prometheus.GaugeValue, float64(m.Protocol.Exported), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].filterCountDesc, prometheus.GaugeValue, float64(m.Protocol.Filtered), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].uptimeDesc, prometheus.GaugeValue, float64(m.Protocol.Uptime), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].runningDesc, prometheus.GaugeValue, m.Protocol.Attributes["running"], m.Protocol.Name)
}
