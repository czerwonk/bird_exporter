package metrics

import (
	"testing"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/stretchr/testify/assert"
)

func TestLabelNames(t *testing.T) {
	s := NewDefaultLabelStrategy(true, `(\w+\s*)=(\s*\w+)`)
	labels := s.LabelNames(&protocol.Protocol{
		Name:         "test",
		Description:  " foo = bar x: y",
		ImportFilter: "in",
		ExportFilter: "out",
		IPVersion:    "6",
		Proto:        protocol.BGP,
	})

	expected := []string{"name", "proto", "ip_version", "import_filter", "export_filter", "foo"}
	assert.Equal(t, expected, labels)
}

func TestLabelValues(t *testing.T) {
	s := NewDefaultLabelStrategy(true, `(\w+\s*)=(\s*\w+)`)
	values := s.LabelValues(&protocol.Protocol{
		Name:         "test",
		Description:  " foo = bar x: y",
		ImportFilter: "in",
		ExportFilter: "out",
		IPVersion:    "6",
		Proto:        protocol.BGP,
	})

	expected := []string{"test", "BGP", "6", "in", "out", "bar"}
	assert.Equal(t, expected, values)
}
