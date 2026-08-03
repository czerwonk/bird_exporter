package metrics

import (
	"testing"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestValidateDescriptionLabelsRegex(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantError  bool
	}{
		{name: "valid", expression: `(\w+)=(\w+)`},
		{name: "invalid syntax", expression: `(`, wantError: true},
		{name: "one capture", expression: `(\w+)=\w+`, wantError: true},
		{name: "three captures", expression: `(\w+)=(\w+)-(\w+)`, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescriptionLabelsRegex(true, tt.expression)
			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestInvalidDescriptionLabelNamesAreDropped(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{name: "invalid name", expression: `([^=]+)=([^ ]+)`, description: `bad-label=value`},
		{name: "reserved name", expression: `(\w+)=(\w+)`, description: `state=value`},
		{name: "duplicate name", expression: `(\w+)=(\w+)`, description: `site=ams site=fra`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewDefaultLabelStrategy(true, tt.expression)
			p := &protocol.Protocol{Name: "peer", Proto: protocol.BGP, Description: tt.description}
			assert.Equal(t, []string{"name", "proto", "ip_version", "import_filter", "export_filter"}, strategy.LabelNames(p))
			assert.Len(t, strategy.LabelValues(p), 5)
		})
	}
}

func TestInvalidRegexDoesNotPanic(t *testing.T) {
	strategy := NewDefaultLabelStrategy(true, `(`)
	p := &protocol.Protocol{Name: "peer", Proto: protocol.BGP, Description: "site=ams"}
	assert.Len(t, strategy.LabelNames(p), 5)
	assert.Len(t, strategy.LabelValues(p), 5)
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
