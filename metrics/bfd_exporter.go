package metrics

import (
	"github.com/czerwonk/bird_exporter/client"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	bfdUpDesc       *prometheus.Desc
	bfdUptimeDesc   *prometheus.Desc
	bfdIntervalDesc *prometheus.Desc
	bfdTimoutDesc   *prometheus.Desc
)

func init() {
	l := []string{"name", "ip", "interface"}
	prefix := "bird_bfd_session_"
	bfdUpDesc = prometheus.NewDesc(prefix+"up", "Session is up", l, nil)
	bfdUptimeDesc = prometheus.NewDesc(prefix+"uptime_seconds", "Session uptime in seconds", l, nil)
	bfdIntervalDesc = prometheus.NewDesc(prefix+"interval_seconds", "Session uptime in seconds", l, nil)
	bfdTimoutDesc = prometheus.NewDesc(prefix+"timeout_seconds", "Session timeout in seconds", l, nil)
}

type bfdMetricExporter struct {
	client client.Client
}

// NewBFDExporter creates a new MetricExporter for BFD metrics
func NewBFDExporter(client client.Client) MetricExporter {
	return &bfdMetricExporter{client: client}
}

func (m *bfdMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- bfdUpDesc
	ch <- bfdUptimeDesc
	ch <- bfdIntervalDesc
	ch <- bfdTimoutDesc
}

func (m *bfdMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newFormat bool) {
	if p.Proto != protocol.BFD {
		return
	}

	sessions, err := m.client.GetBFDSessions(p)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, s := range sessions {
		m.exportSession(s, p.Name, ch)
	}
}

func (m *bfdMetricExporter) exportSession(s *protocol.BFDSession, protocolName string, ch chan<- prometheus.Metric) {
	l := []string{protocolName, s.IP, s.Interface}

	var up float64
	var uptime float64

	if s.Up {
		up = 1
		uptime = float64(s.Since)
	}

	ch <- prometheus.MustNewConstMetric(bfdUpDesc, prometheus.GaugeValue, up, l...)
	ch <- prometheus.MustNewConstMetric(bfdUptimeDesc, prometheus.GaugeValue, uptime, l...)
	ch <- prometheus.MustNewConstMetric(bfdIntervalDesc, prometheus.GaugeValue, s.Interval, l...)
	ch <- prometheus.MustNewConstMetric(bfdTimoutDesc, prometheus.GaugeValue, s.Timeout, l...)
}
