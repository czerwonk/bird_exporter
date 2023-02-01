package main

import (
	"github.com/czerwonk/bird_exporter/client"
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type MetricCollector struct {
	exporters        map[protocol.Proto][]metrics.MetricExporter
	client           *client.BirdClient
	enabledProtocols protocol.Proto
	newFormat        bool
}

func NewMetricCollector(newFormat bool, enabledProtocols protocol.Proto, descriptionLabels bool) *MetricCollector {
	c := getClient()
	var e map[protocol.Proto][]metrics.MetricExporter

	if newFormat {
		e = exportersForDefault(c, descriptionLabels)
	} else {
		e = exportersForLegacy(c)
	}

	return &MetricCollector{
		exporters:        e,
		client:           c,
		enabledProtocols: enabledProtocols,
		newFormat:        newFormat,
	}
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

func exportersForLegacy(c *client.BirdClient) map[protocol.Proto][]metrics.MetricExporter {
	l := metrics.NewLegacyLabelStrategy()

	return map[protocol.Proto][]metrics.MetricExporter{
		protocol.BGP:    {metrics.NewLegacyMetricExporter("bgp4_session", "bgp6_session", l)},
		protocol.Direct: {metrics.NewLegacyMetricExporter("direct4", "direct6", l)},
		protocol.Kernel: {metrics.NewLegacyMetricExporter("kernel4", "kernel6", l)},
		protocol.OSPF:   {metrics.NewLegacyMetricExporter("ospf", "ospfv3", l), metrics.NewOSPFExporter("", c)},
		protocol.Static: {metrics.NewLegacyMetricExporter("static4", "static6", l)},
		protocol.Babel:  {metrics.NewLegacyMetricExporter("babel4", "babel6", l)},
		protocol.RPKI:   {metrics.NewLegacyMetricExporter("rpki4", "rpki6", l)},
		protocol.BFD:    {metrics.NewBFDExporter(c)},
	}
}

func exportersForDefault(c *client.BirdClient, descriptionLabels bool) map[protocol.Proto][]metrics.MetricExporter {
	l := metrics.NewDefaultLabelStrategy(descriptionLabels, *descriptionLabelsRegex)
	e := metrics.NewGenericProtocolMetricExporter("bird_protocol", true, l)

	return map[protocol.Proto][]metrics.MetricExporter{
		protocol.BGP:    {e},
		protocol.Direct: {e},
		protocol.Kernel: {e},
		protocol.OSPF:   {e, metrics.NewOSPFExporter("bird_", c)},
		protocol.Static: {e},
		protocol.Babel:  {e},
		protocol.RPKI:   {e},
		protocol.BFD:    {metrics.NewBFDExporter(c)},
	}
}

var socketQueryDesc = prometheus.NewDesc(
	"bird_socket_query_success",
	"Result of querying bird socket: 0 = failed, 1 = suceeded",
	nil,
	nil,
)

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- socketQueryDesc

	for _, v := range m.exporters {
		for _, e := range v {
			e.Describe(ch)
		}
	}
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {

	protocols, err := m.client.GetProtocols()

	var queryResult float64 = 1
	if err != nil {
		queryResult = 0
	}
	ch <- prometheus.MustNewConstMetric(socketQueryDesc, prometheus.GaugeValue, queryResult)

	if err != nil {
		log.Errorln(err)
		return
	}

	for _, p := range protocols {
		if p.Proto == protocol.PROTO_UNKNOWN || (m.enabledProtocols&p.Proto != p.Proto) {
			continue
		}

		for _, e := range m.exporters[p.Proto] {
			e.Export(p, ch, m.newFormat)
		}
	}
}
