package metrics

import "github.com/czerwonk/bird_exporter/protocol"

type LabelStrategy interface {
	labelNames() []string
	labelValues(p *protocol.Protocol) []string
}
