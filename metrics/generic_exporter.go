package metrics

import (
	"strconv"
	"strings"
	"sync"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

// maxCachedDescSets bounds the desc cache. Description labels can yield a new
// label set per protocol description, so the cache is reset when full.
const maxCachedDescSets = 256

var (
	descSetCacheMu sync.Mutex
	descSetCache   = map[string]*genericDescs{}
)

// GenericProtocolMetricExporter exports metrics retrieved from Bird routing daemon
type GenericProtocolMetricExporter struct {
	labelStrategy LabelStrategy
	prefix        string
}

type genericDescs struct {
	up                          *prometheus.Desc
	importCount                 *prometheus.Desc
	exportCount                 *prometheus.Desc
	filterCount                 *prometheus.Desc
	preferredCount              *prometheus.Desc
	uptime                      *prometheus.Desc
	updatesImportIgnoreCount    *prometheus.Desc
	updatesImportAcceptCount    *prometheus.Desc
	updatesImportFilterCount    *prometheus.Desc
	updatesImportRejectCount    *prometheus.Desc
	updatesImportReceiveCount   *prometheus.Desc
	withdrawsImportIgnoreCount  *prometheus.Desc
	withdrawsImportAcceptCount  *prometheus.Desc
	withdrawsImportFilterCount  *prometheus.Desc
	withdrawsImportRejectCount  *prometheus.Desc
	withdrawsImportReceiveCount *prometheus.Desc
	updatesExportIgnoreCount    *prometheus.Desc
	updatesExportAcceptCount    *prometheus.Desc
	updatesExportFilterCount    *prometheus.Desc
	updatesExportRejectCount    *prometheus.Desc
	updatesExportReceiveCount   *prometheus.Desc
	withdrawsExportIgnoreCount  *prometheus.Desc
	withdrawsExportAcceptCount  *prometheus.Desc
	withdrawsExportFilterCount  *prometheus.Desc
	withdrawsExportRejectCount  *prometheus.Desc
	withdrawsExportReceiveCount *prometheus.Desc
	updatesImportRxLimitCount   *prometheus.Desc
	updatesImportLimitCount     *prometheus.Desc
	updatesExportRxLimitCount   *prometheus.Desc
	updatesExportLimitCount     *prometheus.Desc
	withdrawsImportRxLimitCount *prometheus.Desc
	withdrawsImportLimitCount   *prometheus.Desc
	withdrawsExportRxLimitCount *prometheus.Desc
	withdrawsExportLimitCount   *prometheus.Desc
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

// descsFor returns the descs for a label set, building them only on first use.
func descsFor(prefix string, newNaming bool, labels []string) *genericDescs {
	key := prefix + "\xff" + strconv.FormatBool(newNaming) + "\xff" + strings.Join(labels, "\xff")

	descSetCacheMu.Lock()
	defer descSetCacheMu.Unlock()

	if d, ok := descSetCache[key]; ok {
		return d
	}

	if len(descSetCache) >= maxCachedDescSets {
		descSetCache = map[string]*genericDescs{}
	}

	d := newGenericDescs(prefix, newNaming, labels)
	descSetCache[key] = d
	return d
}

func newGenericDescs(prefix string, newNaming bool, labels []string) *genericDescs {
	d := &genericDescs{}

	d.up = prometheus.NewDesc(prefix+"_up", "Protocol is up", append(labels[:len(labels):len(labels)], "state"), nil)

	if newNaming {
		d.importCount = prometheus.NewDesc(prefix+"_prefix_import_count", "Number of imported routes", labels, nil)
		d.exportCount = prometheus.NewDesc(prefix+"_prefix_export_count", "Number of exported routes", labels, nil)
		d.filterCount = prometheus.NewDesc(prefix+"_prefix_filter_count", "Number of filtered routes", labels, nil)
		d.preferredCount = prometheus.NewDesc(prefix+"_prefix_preferred_count", "Number of preferred routes", labels, nil)
	} else {
		d.importCount = prometheus.NewDesc(prefix+"_prefix_count_import", "Number of imported routes", labels, nil)
		d.exportCount = prometheus.NewDesc(prefix+"_prefix_count_export", "Number of exported routes", labels, nil)
		d.filterCount = prometheus.NewDesc(prefix+"_prefix_count_filter", "Number of filtered routes", labels, nil)
		d.preferredCount = prometheus.NewDesc(prefix+"_prefix_count_preferred", "Number of preferred routes", labels, nil)
	}

	d.uptime = prometheus.NewDesc(prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
	d.updatesImportIgnoreCount = prometheus.NewDesc(prefix+"_changes_update_import_ignore_count", "Number of incoming updates being ignored", labels, nil)
	d.updatesImportAcceptCount = prometheus.NewDesc(prefix+"_changes_update_import_accept_count", "Number of incoming updates being accepted", labels, nil)
	d.updatesImportFilterCount = prometheus.NewDesc(prefix+"_changes_update_import_filter_count", "Number of incoming updates being filtered", labels, nil)
	d.updatesImportRejectCount = prometheus.NewDesc(prefix+"_changes_update_import_reject_count", "Number of incoming updates being rejected", labels, nil)
	d.updatesImportReceiveCount = prometheus.NewDesc(prefix+"_changes_update_import_receive_count", "Number of received updates", labels, nil)
	d.withdrawsImportIgnoreCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_ignore_count", "Number of incoming withdraws being ignored", labels, nil)
	d.withdrawsImportAcceptCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_accept_count", "Number of incoming withdraws being accepted", labels, nil)
	d.withdrawsImportFilterCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_filter_count", "Number of incoming withdraws being filtered", labels, nil)
	d.withdrawsImportRejectCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_reject_count", "Number of incoming withdraws being rejected", labels, nil)
	d.withdrawsImportReceiveCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_receive_count", "Number of received withdraws", labels, nil)
	d.updatesExportIgnoreCount = prometheus.NewDesc(prefix+"_changes_update_export_ignore_count", "Number of outgoing updates being ignored", labels, nil)
	d.updatesExportAcceptCount = prometheus.NewDesc(prefix+"_changes_update_export_accept_count", "Number of outgoing updates being accepted", labels, nil)
	d.updatesExportFilterCount = prometheus.NewDesc(prefix+"_changes_update_export_filter_count", "Number of outgoing updates being filtered", labels, nil)
	d.updatesExportRejectCount = prometheus.NewDesc(prefix+"_changes_update_export_reject_count", "Number of outgoing updates being rejected", labels, nil)
	d.updatesExportReceiveCount = prometheus.NewDesc(prefix+"_changes_update_export_receive_count", "Number of sent updates", labels, nil)
	d.withdrawsExportIgnoreCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_ignore_count", "Number of outgoing withdraws being ignored", labels, nil)
	d.withdrawsExportAcceptCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_accept_count", "Number of outgoing withdraws being accepted", labels, nil)
	d.withdrawsExportFilterCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_filter_count", "Number of outgoing withdraws being filtered", labels, nil)
	d.withdrawsExportRejectCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_reject_count", "Number of outgoing withdraws being rejected", labels, nil)
	d.withdrawsExportReceiveCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_receive_count", "Number of outgoing withdraws", labels, nil)

	d.updatesImportRxLimitCount = prometheus.NewDesc(prefix+"_changes_update_import_rx_limit_count", "Number of incoming updates reaching RX limit", labels, nil)
	d.updatesImportLimitCount = prometheus.NewDesc(prefix+"_changes_update_import_limit_count", "Number of incoming updates reaching limit", labels, nil)
	d.updatesExportRxLimitCount = prometheus.NewDesc(prefix+"_changes_update_export_rx_limit_count", "Number of outgoing updates reaching RX limit", labels, nil)
	d.updatesExportLimitCount = prometheus.NewDesc(prefix+"_changes_update_export_limit_count", "Number of outgoing updates reaching limit", labels, nil)
	d.withdrawsImportRxLimitCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_rx_limit_count", "Number of incoming withdraws reaching RX limit", labels, nil)
	d.withdrawsImportLimitCount = prometheus.NewDesc(prefix+"_changes_withdraw_import_limit_count", "Number of incoming withdraws reaching limit", labels, nil)
	d.withdrawsExportRxLimitCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_rx_limit_count", "Number of outgoing withdraws reaching RX limit", labels, nil)
	d.withdrawsExportLimitCount = prometheus.NewDesc(prefix+"_changes_withdraw_export_limit_count", "Number of outgoing withdraws reaching limit", labels, nil)

	return d
}

func (m *GenericProtocolMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newNaming bool) {
	d := descsFor(m.prefix, newNaming, m.labelStrategy.LabelNames(p))
	l := m.labelStrategy.LabelValues(p)

	ch <- prometheus.MustNewConstMetric(d.up, prometheus.GaugeValue, float64(p.Up), append(l, p.State)...)
	ch <- prometheus.MustNewConstMetric(d.importCount, prometheus.GaugeValue, float64(p.Imported), l...)
	ch <- prometheus.MustNewConstMetric(d.exportCount, prometheus.GaugeValue, float64(p.Exported), l...)
	ch <- prometheus.MustNewConstMetric(d.filterCount, prometheus.GaugeValue, float64(p.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(d.preferredCount, prometheus.GaugeValue, float64(p.Preferred), l...)
	ch <- prometheus.MustNewConstMetric(d.uptime, prometheus.GaugeValue, float64(p.Uptime), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesImportReceiveCount, prometheus.GaugeValue, float64(p.ImportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesImportRejectCount, prometheus.GaugeValue, float64(p.ImportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesImportFilterCount, prometheus.GaugeValue, float64(p.ImportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesImportAcceptCount, prometheus.GaugeValue, float64(p.ImportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesImportIgnoreCount, prometheus.GaugeValue, float64(p.ImportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesExportReceiveCount, prometheus.GaugeValue, float64(p.ExportUpdates.Received), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesExportRejectCount, prometheus.GaugeValue, float64(p.ExportUpdates.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesExportFilterCount, prometheus.GaugeValue, float64(p.ExportUpdates.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesExportAcceptCount, prometheus.GaugeValue, float64(p.ExportUpdates.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(d.updatesExportIgnoreCount, prometheus.GaugeValue, float64(p.ExportUpdates.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsImportReceiveCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsImportRejectCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsImportFilterCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsImportAcceptCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsImportIgnoreCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Ignored), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsExportReceiveCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Received), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsExportRejectCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Rejected), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsExportFilterCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Filtered), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsExportAcceptCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Accepted), l...)
	ch <- prometheus.MustNewConstMetric(d.withdrawsExportIgnoreCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Ignored), l...)

	if p.RouteChangeFormatV3 {
		ch <- prometheus.MustNewConstMetric(d.updatesImportRxLimitCount, prometheus.GaugeValue, float64(p.ImportUpdates.RxLimit), l...)
		ch <- prometheus.MustNewConstMetric(d.updatesImportLimitCount, prometheus.GaugeValue, float64(p.ImportUpdates.Limit), l...)
		ch <- prometheus.MustNewConstMetric(d.updatesExportRxLimitCount, prometheus.GaugeValue, float64(p.ExportUpdates.RxLimit), l...)
		ch <- prometheus.MustNewConstMetric(d.updatesExportLimitCount, prometheus.GaugeValue, float64(p.ExportUpdates.Limit), l...)
		ch <- prometheus.MustNewConstMetric(d.withdrawsImportRxLimitCount, prometheus.GaugeValue, float64(p.ImportWithdraws.RxLimit), l...)
		ch <- prometheus.MustNewConstMetric(d.withdrawsImportLimitCount, prometheus.GaugeValue, float64(p.ImportWithdraws.Limit), l...)
		ch <- prometheus.MustNewConstMetric(d.withdrawsExportRxLimitCount, prometheus.GaugeValue, float64(p.ExportWithdraws.RxLimit), l...)
		ch <- prometheus.MustNewConstMetric(d.withdrawsExportLimitCount, prometheus.GaugeValue, float64(p.ExportWithdraws.Limit), l...)
	}
}
