package metrics

import "github.com/czerwonk/bird_exporter/protocol"

type LegacyLabelStrategy struct {
}

func (*LegacyLabelStrategy) labelNames() []string {
	return []string{"name"}
}

func (*LegacyLabelStrategy) labelValues(p *protocol.Protocol) []string {
	return []string{p.Name}
}
