package metrics

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

// GenericProtocolMetricExporter exports metrics retrieved from Bird routing daemon
type GenericProtocolMetricExporter struct {
	labelStrategy LabelStrategy
	prefix        string
}

// NewGenericProtocolMetricExporter creates a new instance of GenericProtocolMetricExporter
func NewGenericProtocolMetricExporter(prefix string, newNaming bool, labelStrategy LabelStrategy) *GenericProtocolMetricExporter {
	m := &GenericProtocolMetricExporter{
		prefix:        prefix,
		labelStrategy: labelStrategy,
	}

	return m
}

func (m *GenericProtocolMetricExporter) Describe(ch chan<- *prometheus.Desc) {
}

func (m *GenericProtocolMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newNaming bool) {
	labels := m.labelStrategy.LabelNames(p)

	var importCountDesc *prometheus.Desc
	var exportCountDesc *prometheus.Desc
	var filterCountDesc *prometheus.Desc

	upDesc := prometheus.NewDesc(m.prefix+"_up", "Protocol is up", append(labels, "state"), nil)

	if newNaming {
		importCountDesc = prometheus.NewDesc(m.prefix+"_prefix_import_count", "Number of imported routes", labels, nil)
		exportCountDesc = prometheus.NewDesc(m.prefix+"_prefix_export_count", "Number of exported routes", labels, nil)
		filterCountDesc = prometheus.NewDesc(m.prefix+"_prefix_filter_count", "Number of filtered routes", labels, nil)
	} else {
		importCountDesc = prometheus.NewDesc(m.prefix+"_prefix_count_import", "Number of imported routes", labels, nil)
		exportCountDesc = prometheus.NewDesc(m.prefix+"_prefix_count_export", "Number of exported routes", labels, nil)
		filterCountDesc = prometheus.NewDesc(m.prefix+"_prefix_count_filter", "Number of filtered routes", labels, nil)
	}

	uptimeDesc := prometheus.NewDesc(m.prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
	updatesImportIgnoreCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_import_ignore_count", "Number of incoming updates being ignored", labels, nil)
	updatesImportAcceptCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_import_accept_count", "Number of incoming updates being accepted", labels, nil)
	updatesImportFilterCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_import_filter_count", "Number of incoming updates being filtered", labels, nil)
	updatesImportRejectCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_import_reject_count", "Number of incoming updates being rejected", labels, nil)
	updatesImportReceiveCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_import_receive_count", "Number of received updates", labels, nil)
	withdrawsImportIgnoreCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_import_ignore_count", "Number of incoming withdraws being ignored", labels, nil)
	withdrawsImportAcceptCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_import_accept_count", "Number of incoming withdraws being accepted", labels, nil)
	withdrawsImportFilterCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_import_filter_count", "Number of incoming withdraws being filtered", labels, nil)
	withdrawsImportRejectCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_import_reject_count", "Number of incoming withdraws being rejected", labels, nil)
	withdrawsImportReceiveCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_import_receive_count", "Number of received withdraws", labels, nil)
	updatesExportIgnoreCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_export_ignore_count", "Number of outgoing updates being ignored", labels, nil)
	updatesExportAcceptCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_export_accept_count", "Number of outgoing updates being accepted", labels, nil)
	updatesExportFilterCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_export_filter_count", "Number of outgoing updates being filtered", labels, nil)
	updatesExportRejectCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_export_reject_count", "Number of outgoing updates being rejected", labels, nil)
	updatesExportReceiveCountDesc := prometheus.NewDesc(m.prefix+"_changes_update_export_receive_count", "Number of sent updates", labels, nil)
	withdrawsExportIgnoreCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_export_ignore_count", "Number of outgoing withdraws being ignored", labels, nil)
	withdrawsExportAcceptCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_export_accept_count", "Number of outgoing withdraws being accepted", labels, nil)
	withdrawsExportFilterCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_export_filter_count", "Number of outgoing withdraws being filtered", labels, nil)
	withdrawsExportRejectCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_export_reject_count", "Number of outgoing withdraws being rejected", labels, nil)
	withdrawsExportReceiveCountDesc := prometheus.NewDesc(m.prefix+"_changes_withdraw_export_receive_count", "Number of outgoing withdraws", labels, nil)

	l := m.labelStrategy.LabelValues(p)
	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, float64(p.Up), append(l, p.State)...)
	ch <- prometheus.MustNewConstMetric(importCountDesc, prometheus.GaugeValue, float64(p.Imported), l...)
	ch <- prometheus.MustNewConstMetric(exportCountDesc, prometheus.GaugeValue, float64(p.Exported), l...)
	ch <- prometheus.MustNewConstMetric(filterCountDesc, prometheus.GaugeValue, float64(p.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, float64(p.Uptime), l...)
	ch <- prometheus.MustNewConstMetric(updatesImportReceiveCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(updatesImportRejectCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(updatesImportFilterCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(updatesImportAcceptCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(updatesImportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ImportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(updatesExportReceiveCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(updatesExportRejectCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(updatesExportFilterCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(updatesExportAcceptCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(updatesExportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ExportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsImportReceiveCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsImportRejectCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsImportFilterCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsImportAcceptCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsImportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ImportWithdraws.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsExportReceiveCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsExportRejectCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsExportFilterCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsExportAcceptCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(withdrawsExportIgnoreCountDesc, prometheus.GaugeValue, float64(p.ExportWithdraws.Ignored), l...)
}
