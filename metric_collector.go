package main

import (
	"github.com/czerwonk/bird_exporter/client"
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type MetricCollector struct {
	exporters        map[int][]metrics.MetricExporter
	client           *client.BirdClient
	enabledProtocols int
}

func NewMetricCollector(newFormat bool, enabledProtocols int) *MetricCollector {
	c := getClient()
	var e map[int][]metrics.MetricExporter

	if newFormat {
		e = exportersForDefault(c)
	} else {
		e = exportersForLegacy(c)
	}

	return &MetricCollector{exporters: e, client: c, enabledProtocols: enabledProtocols}
}

func getClient() *client.BirdClient {
	o := &client.BirdClientOptions{
		BirdSocket:   *birdSocket,
		Bird6Socket:  *bird6Socket,
		Bird6Enabled: *bird6Enabled,
		BirdEnabled:  *birdEnabled,
		BirdV2:       *birdV2,
	}

	return &client.BirdClient{Options: o}
}

func exportersForLegacy(c *client.BirdClient) map[int][]metrics.MetricExporter {
	l := &metrics.LegacyLabelStrategy{}

	return map[int][]metrics.MetricExporter{
		protocol.BGP:    {metrics.NewLegacyMetricExporter("bgp4_session", "bgp6_session", l)},
		protocol.Direct: {metrics.NewLegacyMetricExporter("direct4", "direct6", l)},
		protocol.Kernel: {metrics.NewLegacyMetricExporter("kernel4", "kernel6", l)},
		protocol.OSPF:   {metrics.NewLegacyMetricExporter("ospf", "ospfv3", l), metrics.NewOspfExporter("", c)},
		protocol.Static: {metrics.NewLegacyMetricExporter("static4", "static6", l)},
	}
}

func exportersForDefault(c *client.BirdClient) map[int][]metrics.MetricExporter {
	l := &metrics.DefaultLabelStrategy{}
	e := metrics.NewGenericProtocolMetricExporter("bird_protocol", true, l)

	return map[int][]metrics.MetricExporter{
		protocol.BGP:    {e},
		protocol.Direct: {e},
		protocol.Kernel: {e},
		protocol.OSPF:   {e, metrics.NewOspfExporter("bird_", c)},
		protocol.Static: {e},
	}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range m.exporters {
		for _, e := range v {
			e.Describe(ch)
		}
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	protocols, err := m.client.GetProtocols()
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, p := range protocols {
		if p.Proto == protocol.PROTO_UNKNOWN || (m.enabledProtocols&p.Proto != p.Proto) {
			continue
		}

		for _, e := range m.exporters[p.Proto] {
			e.Export(p, ch)
		}
	}
}
