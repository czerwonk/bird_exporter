package metrics

import "github.com/czerwonk/bird_exporter/protocol"

// LabelStrategy abstracts the label generation for protocol metrics
type LabelStrategy interface {
	// LabelNames is the list of label names
	LabelNames() []string
	
	// Label values is the list of values for the labels specified in `LabelNames()` 
	LabelValues(p *protocol.Protocol) []string
}
