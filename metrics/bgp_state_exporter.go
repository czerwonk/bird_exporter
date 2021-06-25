package metrics

import (
	"github.com/czerwonk/bird_exporter/client"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type bgpStateMetricExporter struct {
	prefix string
	client client.Client
}

// NewBGPStateExporter creates a new MetricExporter for BGP metrics
func NewBGPStateExporter(prefix string, client client.Client) MetricExporter {
	return &bgpStateMetricExporter{prefix: prefix, client: client}
}

func (m *bgpStateMetricExporter) Describe(ch chan<- *prometheus.Desc) {
}

func (m *bgpStateMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newFormat bool) {

	labels := []string{"name", "proto", "state"}
	bgpstateDesc := prometheus.NewDesc(m.prefix+"bgp_state_count", "Number of BGP connections at each state", labels, nil)
	state, err := m.client.GetBGPStates(p)
	if err != nil {
		log.Errorln(err)
		return
	}
	if state != nil {
		l := []string{state.Name, "BGP", state.State}
		ch <- prometheus.MustNewConstMetric(bgpstateDesc, prometheus.GaugeValue, float64(1), l...)
	}
}
