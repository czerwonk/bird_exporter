package metrics

import (
	"github.com/czerwonk/bird_exporter/protocol"
)

type DefaultLabelStrategy struct {
}

func (*DefaultLabelStrategy) labelNames() []string {
	return []string{"name", "proto", "ip_version"}
}

func (*DefaultLabelStrategy) labelValues(p *protocol.Protocol) []string {
	return []string{p.Name, protoString(p), p.IpVersion}
}

func protoString(p *protocol.Protocol) string {
	switch p.Proto {
	case protocol.BGP:
		return "BGP"
	case protocol.OSPF:
		if p.IpVersion == "4" {
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
