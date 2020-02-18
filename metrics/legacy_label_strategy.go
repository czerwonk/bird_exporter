package metrics

import "github.com/czerwonk/bird_exporter/protocol"

type LegacyLabelStrategy struct {
}

func NewLegacyLabelStrategy() *LegacyLabelStrategy {
	return &LegacyLabelStrategy{}
}

func (*LegacyLabelStrategy) LabelNames(p *protocol.Protocol) []string {
	return []string{"name"}
}

func (*LegacyLabelStrategy) LabelValues(p *protocol.Protocol) []string {
	return []string{p.Name}
}
