package protocol

import "github.com/prometheus/client_golang/prometheus"

type GenericProtocolMetricExporter struct {
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
	m.updatesImportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_ignore_count", "Number of incoming updates beeing ignored", labels, nil)
	m.updatesImportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_accept_count", "Number of incoming updates beeing accepted", labels, nil)
	m.updatesImportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_filter_count", "Number of incoming updates beeing filtered", labels, nil)
	m.updatesImportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_reject_count", "Number of incoming updates beeing rejected", labels, nil)
	m.updatesImportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_update_import_receive_count", "Number of received updates", labels, nil)
	m.withdrawsImportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_ignore_count", "Number of incoming withdraws beeing ignored", labels, nil)
	m.withdrawsImportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_accept_count", "Number of incoming withdraws beeing accepted", labels, nil)
	m.withdrawsImportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_filter_count", "Number of incoming withdraws beeing filtered", labels, nil)
	m.withdrawsImportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_reject_count", "Number of incoming withdraws beeing rejected", labels, nil)
	m.withdrawsImportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_import_receive_count", "Number of received withdraws", labels, nil)
	m.updatesExportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_ignore_count", "Number of outgoing updates beeing ignored", labels, nil)
	m.updatesExportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_accept_count", "Number of outgoing updates beeing accepted", labels, nil)
	m.updatesExportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_filter_count", "Number of outgoing updates beeing filtered", labels, nil)
	m.updatesExportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_reject_count", "Number of outgoing updates beeing rejected", labels, nil)
	m.updatesExportReceiveCountDesc = prometheus.NewDesc(prefix+"_changes_update_export_receive_count", "Number of sent updates", labels, nil)
	m.withdrawsExportIgnoreCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_ignore_count", "Number of outgoing withdraws beeing ignored", labels, nil)
	m.withdrawsExportAcceptCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_accept_count", "Number of outgoing withdraws beeing accepted", labels, nil)
	m.withdrawsExportFilterCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_filter_count", "Number of outgoing withdraws beeing filtered", labels, nil)
	m.withdrawsExportRejectCountDesc = prometheus.NewDesc(prefix+"_changes_withdraw_export_reject_count", "Number of outgoing withdraws beeing rejected", labels, nil)
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

func (m *GenericProtocolMetricExporter) Export(protocol *Protocol, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.upDesc, prometheus.GaugeValue, float64(protocol.Up), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.importCountDesc, prometheus.GaugeValue, float64(protocol.Imported), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.exportCountDesc, prometheus.GaugeValue, float64(protocol.Exported), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.filterCountDesc, prometheus.GaugeValue, float64(protocol.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.preferredCountDesc, prometheus.GaugeValue, float64(protocol.Preferred), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.uptimeDesc, prometheus.GaugeValue, float64(protocol.Uptime), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesImportReceiveCountDesc, prometheus.GaugeValue, float64(protocol.ImportUpdates.Received), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesImportRejectCountDesc, prometheus.GaugeValue, float64(protocol.ImportUpdates.Rejected), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesImportFilterCountDesc, prometheus.GaugeValue, float64(protocol.ImportUpdates.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesImportAcceptCountDesc, prometheus.GaugeValue, float64(protocol.ImportUpdates.Accepted), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesImportIgnoreCountDesc, prometheus.GaugeValue, float64(protocol.ImportUpdates.Ignored), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesExportReceiveCountDesc, prometheus.GaugeValue, float64(protocol.ExportUpdates.Received), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesExportRejectCountDesc, prometheus.GaugeValue, float64(protocol.ExportUpdates.Rejected), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesExportFilterCountDesc, prometheus.GaugeValue, float64(protocol.ExportUpdates.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesExportAcceptCountDesc, prometheus.GaugeValue, float64(protocol.ExportUpdates.Accepted), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.updatesExportIgnoreCountDesc, prometheus.GaugeValue, float64(protocol.ExportUpdates.Ignored), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportReceiveCountDesc, prometheus.GaugeValue, float64(protocol.ImportWithdraws.Received), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportRejectCountDesc, prometheus.GaugeValue, float64(protocol.ImportWithdraws.Rejected), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportFilterCountDesc, prometheus.GaugeValue, float64(protocol.ImportWithdraws.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportAcceptCountDesc, prometheus.GaugeValue, float64(protocol.ImportWithdraws.Accepted), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsImportIgnoreCountDesc, prometheus.GaugeValue, float64(protocol.ImportWithdraws.Ignored), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportReceiveCountDesc, prometheus.GaugeValue, float64(protocol.ExportWithdraws.Received), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportRejectCountDesc, prometheus.GaugeValue, float64(protocol.ExportWithdraws.Rejected), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportFilterCountDesc, prometheus.GaugeValue, float64(protocol.ExportWithdraws.Filtered), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportAcceptCountDesc, prometheus.GaugeValue, float64(protocol.ExportWithdraws.Accepted), protocol.Name)
	ch <- prometheus.MustNewConstMetric(m.withdrawsExportIgnoreCountDesc, prometheus.GaugeValue, float64(protocol.ExportWithdraws.Ignored), protocol.Name)
}
