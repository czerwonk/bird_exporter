package metrics

import "github.com/czerwonk/bird_exporter/protocol"

type LabelStrategy interface {
	LabelNames() []string
	LabelValues(p *protocol.Protocol) []string
}
