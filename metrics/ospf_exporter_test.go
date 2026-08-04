package metrics

import (
	"testing"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestOSPFExporterRejectsUnsupportedIPVersionWithoutPanic(t *testing.T) {
	exporter := NewOSPFExporter("bird_", nil)
	metricChannel := make(chan prometheus.Metric, 1)

	require.NotPanics(t, func() {
		exporter.Export(&protocol.Protocol{Name: "ospf_unknown", Proto: protocol.OSPF}, metricChannel, true)
	})

	metric := <-metricChannel
	require.Error(t, metric.Write(&dto.Metric{}))
}
