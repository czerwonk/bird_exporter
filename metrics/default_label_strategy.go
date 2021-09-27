package metrics

import (
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

// DefaultLabelStrategy defines the labels to add to an metric and its data retrieval method
type DefaultLabelStrategy struct {
	descriptionLabels bool
}

func NewDefaultLabelStrategy(descriptionLabels bool) *DefaultLabelStrategy {
	return &DefaultLabelStrategy{
		descriptionLabels: descriptionLabels,
	}
}

// LabelNames returns the list of label names
func (d *DefaultLabelStrategy) LabelNames(p *protocol.Protocol) []string {
	res := []string{"name", "proto", "ip_version", "import_filter", "export_filter"}
	if d.descriptionLabels && p.Description != "" {
		res = append(res, labelKeysFromDescription(p.Description)...)
	}

	return res
}

// LabelValues returns the values for a protocol
func (d *DefaultLabelStrategy) LabelValues(p *protocol.Protocol) []string {
	res := []string{p.Name, protoString(p), p.IPVersion, p.ImportFilter, p.ExportFilter}
	if d.descriptionLabels && p.Description != "" {
		res = append(res, labelValuesFromDescription(p.Description)...)
	}

	return res
}

func labelKeysFromDescription(desc string) (res []string) {
	for _, x := range strings.Split(desc, ",") {
		tmp := strings.Split(x, "=")
		if len(tmp) != 2 {
			continue
		}

		res = append(res, tmp[0])
	}

	return
}

func labelValuesFromDescription(desc string) (res []string) {
	for _, x := range strings.Split(desc, ",") {
		tmp := strings.Split(x, "=")
		if len(tmp) != 2 {
			continue
		}

		res = append(res, tmp[1])
	}

	return
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
	case protocol.Babel:
		return "Babel"
	}

	return ""
}
