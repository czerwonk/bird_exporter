package metrics

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type GenericProtocolMetricExporter struct {
	labelStrategy                   LabelStrategy
	upDesc                          *prometheus.Desc
	importCountDesc                 *prometheus.Desc
	exportCountDesc                 *prometheus.Desc
	filterCountDesc                 *prometheus.Desc
	preferredCountDesc              *prometheus.Desc
	uptimeDesc                      *prometheus.Desc
	updatesImportReceiveCountDesc   *prometheus.Desc
	updatesImportRejectCountDesc    *prometheus.Desc
	updatesImportFilterCountDesc    *prometheus.Desc
	updatesImportIgnoreCountDesc    *prometheus.Desc
	updatesImportAcceptCountDesc    *prometheus.Desc
	withdrawsImportReceiveCountDesc *prometheus.Desc
	withdrawsImportRejectCountDesc  *prometheus.Desc
	withdrawsImportFilterCountDesc  *prometheus.Desc
	withdrawsImportIgnoreCountDesc  *prometheus.Desc
	withdrawsImportAcceptCountDesc  *prometheus.Desc
	updatesExportReceiveCountDesc   *prometheus.Desc
	updatesExportRejectCountDesc    *prometheus.Desc
	updatesExportFilterCountDesc    *prometheus.Desc
	updatesExportIgnoreCountDesc    *prometheus.Desc
	updatesExportAcceptCountDesc    *prometheus.Desc
	withdrawsExportReceiveCountDesc *prometheus.Desc
	withdrawsExportRejectCountDesc  *prometheus.Desc
	withdrawsExportFilterCountDesc  *prometheus.Desc
	withdrawsExportIgnoreCountDesc  *prometheus.Desc
	withdrawsExportAcceptCountDesc  *prometheus.Desc
}

func NewGenericProtocolMetricExporter(prefix string, newNaming bool, labelStrategy LabelStrategy) *GenericProtocolMetricExporter {
	m := &GenericProtocolMetricExporter{labelStrategy: labelStrategy}
	m.initDesc(prefix, newNaming)

	return m
}

func (m *GenericProtocolMetricExporter) initDesc(prefix string, newNaming bool) {
	labels := m.labelStrategy.LabelNames()
	m.upDesc = prometheus.NewDesc(prefix+"_up", "Protocol is up", labels, nil)

	if newNaming {
		m.importCountDesc = prometheus.NewDesc(prefix+"_prefix_import_count", "Number of imported routes", labels, nil)
		m.exportCountDesc = prometheus.NewDesc(prefix+"_prefix_export_count", "Number of exported routes", labels, nil)
		m.filterCountDesc = prometheus.NewDesc(prefix+"_prefix_filter_count", "Number of filtered routes", labels, nil)
		m.preferredCountDesc = prometheus.NewDesc(prefix+"_prefix_preferred_count", "Number of preferred routes", labels, nil)
	} else {
		m.importCountDesc = prometheus.NewDesc(prefix+"_prefix_count_import", "Number of imported routes", labels, nil)
		m.exportCountDesc = prometheus.NewDesc(prefix+"_prefix_count_export", "Number of exported routes", labels, nil)
		m.filterCountDesc = prometheus.NewDesc(prefix+"_prefix_count_filter", "Number of filtered routes", labels, nil)
		m.preferredCountDesc = prometheus.NewDesc(prefix+"_prefix_count_preferred", "Number of preferred routes", labels, nil)
	}

	m.uptimeDesc = prometheus.NewDesc(prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
	m.updatesImportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_ignore_count", "Number of incoming updates being ignored", labels, nil)
	m.updatesImportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_accept_count", "Number of incoming updates being accepted", labels, nil)
	m.updatesImportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_filter_count", "Number of incoming updates being filtered", labels, nil)
	m.updatesImportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_reject_count", "Number of incoming updates being rejected", labels, nil)
	m.updatesImportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_receive_count", "Number of received updates", labels, nil)
	m.withdrawsImportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_ignore_count", "Number of incoming withdraws being ignored", labels, nil)
	m.withdrawsImportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_accept_count", "Number of incoming withdraws being accepted", labels, nil)
	m.withdrawsImportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_filter_count", "Number of incoming withdraws being filtered", labels, nil)
	m.withdrawsImportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_reject_count", "Number of incoming withdraws being rejected", labels, nil)
	m.withdrawsImportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_receive_count", "Number of received withdraws", labels, nil)
	m.updatesExportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_ignore_count", "Number of outgoing updates being ignored", labels, nil)
	m.updatesExportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_accept_count", "Number of outgoing updates being accepted", labels, nil)
	m.updatesExportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_filter_count", "Number of outgoing updates being filtered", labels, nil)
	m.updatesExportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_reject_count", "Number of outgoing updates being rejected", labels, nil)
	m.updatesExportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_receive_count", "Number of sent updates", labels, nil)
	m.withdrawsExportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_ignore_count", "Number of outgoing withdraws being ignored", labels, nil)
	m.withdrawsExportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_accept_count", "Number of outgoing withdraws being accepted", labels, nil)
	m.withdrawsExportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_filter_count", "Number of outgoing withdraws being filtered", labels, nil)
	m.withdrawsExportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_reject_count", "Number of outgoing withdraws being rejected", labels, nil)
	m.withdrawsExportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_receive_count", "Number of outgoing withdraws", labels, nil)
}

func (m *GenericProtocolMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.upDesc
	ch <- m.importCountDesc
	ch <- m.exportCountDesc
	ch <- m.filterCountDesc
	ch <- m.preferredCountDesc
	ch <- m.uptimeDesc
	ch <- m.updatesImportReceiveCountDesc
	ch <- m.updatesImportRejectCountDesc
	ch <- m.updatesImportFilterCountDesc
	ch <- m.updatesImportIgnoreCountDesc
	ch <- m.updatesImportAcceptCountDesc
	ch <- m.updatesExportReceiveCountDesc
	ch <- m.updatesExportRejectCountDesc
	ch <- m.updatesExportFilterCountDesc
	ch <- m.updatesExportIgnoreCountDesc
	ch <- m.updatesExportAcceptCountDesc
	ch <- m.withdrawsImportIgnoreCountDesc
	ch <- m.withdrawsImportAcceptCountDesc
	ch <- m.withdrawsImportFilterCountDesc
	ch <- m.withdrawsImportRejectCountDesc
	ch <- m.withdrawsImportReceiveCountDesc
	ch <- m.withdrawsExportIgnoreCountDesc
	ch <- m.withdrawsExportAcceptCountDesc
	ch <- m.withdrawsExportFilterCountDesc
	ch <- m.withdrawsExportRejectCountDesc
	ch <- m.withdrawsExportReceiveCountDesc
}

func (m *GenericProtocolMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric) {
	l := m.labelStrategy.LabelValues(p)
	ch <- prometheus.MustNewConstMetric(m.upDesc, prometheus.GaugeValue, float64(p.Up), l...)
	ch <- prometheus.MustNewConstMetric(m.importCountDesc, prometheus.GaugeValue, float64(p.Imported), l...)
	ch <- prometheus.MustNewConstMetric(m.exportCountDesc, prometheus.GaugeValue, float64(p.Exported), l...)
	ch <- prometheus.MustNewConstMetric(m.filterCountDesc, prometheus.GaugeValue, float64(p.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(m.preferredCountDesc, prometheus.GaugeValue, float64(p.Preferred), l...)
	ch <- prometheus.MustNewConstMetric(m.uptimeDesc, prometheus.GaugeValue, float64(p.Uptime), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesImportReceiveCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesImportRejectCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesImportFilterCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesImportAcceptCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesImportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesExportReceiveCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesExportRejectCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesExportFilterCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesExportAcceptCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(m.updatesExportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportReceiveCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportRejectCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportFilterCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportAcceptCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportReceiveCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportRejectCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportFilterCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportAcceptCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Ignored), l...)
}
