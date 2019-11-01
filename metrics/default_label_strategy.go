package metrics

import (
	"github.com/czerwonk/bird_exporter/protocol"
)

// DefaultLabelStrategy defines the labels to add to an metric and its data retrieval method
type DefaultLabelStrategy struct {
}

// LabelNames returns the list of label names
func (*DefaultLabelStrategy) LabelNames() []string {
	return []string{"name", "proto", "ip_version", "import_filter", "export_filter"}
}

// LabelValues returns the values for a protocol
func (*DefaultLabelStrategy) LabelValues(p *protocol.Protocol) []string {
	return []string{p.Name, protoString(p), p.IPVersion, p.ImportFilter, p.ExportFilter}
}

func protoString(p *protocol.Protocol) string {
	switch p.Proto {
	case protocol.BGP:
		return "BGP"
	case protocol.OSPF:
		if p.IPVersion == "4" {
			return "OSPF"
		}
		return "OSPFv3"
	case protocol.Static:
		return "Static"
	case protocol.Kernel:
		return "Kernel"
	case protocol.Direct:
		return "Direct"
	}

	return ""
}
