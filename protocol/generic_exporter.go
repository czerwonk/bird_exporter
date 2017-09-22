package protocol

import "github.com/prometheus/client_golang/prometheus"

type GenericProtocolMetricExporter struct {
	upDesc             *prometheus.Desc
	importCountDesc    *prometheus.Desc
	exportCountDesc    *prometheus.Desc
	filterCountDesc    *prometheus.Desc
	preferredCountDesc *prometheus.Desc
	uptimeDesc         *prometheus.Desc
}

func NewGenericProtocolMetricExporter(prefix string) *GenericProtocolMetricExporter {
	m := &GenericProtocolMetricExporter{}
	m.initDesc(prefix)

	return m
}

func (m *GenericProtocolMetricExporter) initDesc(prefix string) {
	labels := []string{"name"}
	m.upDesc = prometheus.NewDesc(prefix+"_up", "Protocol is up", labels, nil)
	m.importCountDesc = prometheus.NewDesc(prefix+"_prefix_count_import", "Number of imported routes", labels, nil)
	m.exportCountDesc = prometheus.NewDesc(prefix+"_prefix_count_export", "Number of exported routes", labels, nil)
	m.filterCountDesc = prometheus.NewDesc(prefix+"_prefix_count_filter", "Number of filtered routes", labels, nil)
	m.preferredCountDesc = prometheus.NewDesc(prefix+"_prefix_count_preferred", "Number of preferred routes", labels, nil)
	m.uptimeDesc = prometheus.NewDesc(prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
}

func (m *GenericProtocolMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.upDesc
	ch <- m.importCountDesc
	ch <- m.exportCountDesc
	ch <- m.filterCountDesc
	ch <- m.preferredCountDesc
	ch <- m.uptimeDesc
}

func (m *GenericProtocolMetricExporter) Export(protocol *Protocol, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.upDesc, prometheus.GaugeValue, float64(protocol.Up), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.importCountDesc, prometheus.GaugeValue, float64(protocol.Imported), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.exportCountDesc, prometheus.GaugeValue, float64(protocol.Exported), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.filterCountDesc, prometheus.GaugeValue, float64(protocol.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.preferredCountDesc, prometheus.GaugeValue, float64(protocol.Preferred), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.uptimeDesc, prometheus.GaugeValue, float64(protocol.Uptime), protocol.Name)
}
