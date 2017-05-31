package bgp

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
}

type BgpMetric struct {
	Protocol *protocol.Protocol
}

func init() {
	descriptions = make(map[int]*desc)
	descriptions[4] = getDesc("bgp4")
	descriptions[6] = getDesc("bgp6")
}

func getDesc(prefix string) *desc {
	labels := []string{"name"}

	d := &desc{}
	d.upDesc = prometheus.NewDesc(prefix+"_session_up", "Protocol is up", labels, nil)
	d.importCountDesc = prometheus.NewDesc(prefix+"_session_prefix_count_import", "Number of imported routes", labels, nil)
	d.exportCountDesc = prometheus.NewDesc(prefix+"_session_prefix_count_export", "Number of exported routes", labels, nil)
	d.filterCountDesc = prometheus.NewDesc(prefix+"_session_prefix_count_filter", "Number of filtered routes", labels, nil)
	d.uptimeDesc = prometheus.NewDesc(prefix+"_session_uptime", "Uptime of the protocol in seconds", labels, nil)

	return d
}

func (m *BgpMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- descriptions[m.Protocol.IpVersion].upDesc
	ch <- descriptions[m.Protocol.IpVersion].importCountDesc
	ch <- descriptions[m.Protocol.IpVersion].exportCountDesc
	ch <- descriptions[m.Protocol.IpVersion].filterCountDesc
	ch <- descriptions[m.Protocol.IpVersion].uptimeDesc
}

func (m *BgpMetric) GetMetrics(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].upDesc, prometheus.GaugeValue, float64(m.Protocol.Up), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].importCountDesc, prometheus.GaugeValue, float64(m.Protocol.Imported), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].exportCountDesc, prometheus.GaugeValue, float64(m.Protocol.Exported), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].filterCountDesc, prometheus.GaugeValue, float64(m.Protocol.Filtered), m.Protocol.Name)
	ch <- prometheus.MustNewConstMetric(descriptions[m.Protocol.IpVersion].uptimeDesc, prometheus.GaugeValue, float64(m.Protocol.Uptime), m.Protocol.Name)
}
