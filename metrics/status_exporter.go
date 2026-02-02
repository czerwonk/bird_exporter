package metrics

import (
	"github.com/czerwonk/bird_exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type StatusExporter struct {
	client     *client.BirdClient
	socketPath string

	upDesc           *prometheus.Desc
	infoDesc         *prometheus.Desc
	lastRebootDesc   *prometheus.Desc
	lastReconfigDesc *prometheus.Desc
	serverTimeDesc   *prometheus.Desc
}

func NewStatusExporter(c *client.BirdClient, socketPath string) *StatusExporter {
	return &StatusExporter{
		client:     c,
		socketPath: socketPath,

		upDesc: prometheus.NewDesc(
			"bird_daemon_up",
			"Whether the BIRD daemon is up (1) or down (0)",
			nil, nil,
		),
		infoDesc: prometheus.NewDesc(
			"bird_daemon_info",
			"Static information about BIRD",
			[]string{"router_id", "version"}, nil,
		),
		lastRebootDesc: prometheus.NewDesc(
			"bird_last_reboot_timestamp_seconds",
			"Timestamp of the last BIRD reboot",
			nil, nil,
		),
		lastReconfigDesc: prometheus.NewDesc(
			"bird_last_reconfig_timestamp_seconds",
			"Timestamp of the last BIRD reconfiguration",
			nil, nil,
		),
		serverTimeDesc: prometheus.NewDesc(
			"bird_server_time_timestamp_seconds",
			"Server time reported by BIRD",
			nil, nil,
		),
	}
}

func (e *StatusExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.upDesc
	ch <- e.infoDesc
	ch <- e.lastRebootDesc
	ch <- e.lastReconfigDesc
	ch <- e.serverTimeDesc
}

func (e *StatusExporter) Collect(ch chan<- prometheus.Metric) {
	s, err := e.client.StatusFromSocket(e.socketPath)
	if err != nil || s == nil {
		ch <- prometheus.MustNewConstMetric(e.upDesc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.upDesc, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(e.infoDesc, prometheus.GaugeValue, 1, s.RouterID, s.Version)

	if !s.LastReboot.IsZero() {
		ch <- prometheus.MustNewConstMetric(e.lastRebootDesc, prometheus.GaugeValue, float64(s.LastReboot.Unix()))
	}
	if !s.LastReconfig.IsZero() {
		ch <- prometheus.MustNewConstMetric(e.lastReconfigDesc, prometheus.GaugeValue, float64(s.LastReconfig.Unix()))
	}
	if !s.ServerTime.IsZero() {
		ch <- prometheus.MustNewConstMetric(e.serverTimeDesc, prometheus.GaugeValue, float64(s.ServerTime.Unix()))
	}
}
