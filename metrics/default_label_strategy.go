package metrics

import (
	"regexp"

	"github.com/czerwonk/bird_exporter/protocol"
)

// DefaultLabelStrategy defines the labels to add to an metric and its data retrieval method
type DefaultLabelStrategy struct {
	descriptionLabels      bool
	descriptionLabelsRegex string
}

func NewDefaultLabelStrategy(descriptionLabels bool, descriptionLabelsRegex string) *DefaultLabelStrategy {
	return &DefaultLabelStrategy{
		descriptionLabels:      descriptionLabels,
		descriptionLabelsRegex: descriptionLabelsRegex,
	}
}

// LabelNames returns the list of label names
func (d *DefaultLabelStrategy) LabelNames(p *protocol.Protocol) []string {
	res := []string{"name", "proto", "ip_version", "import_filter", "export_filter"}
	if d.descriptionLabels && p.Description != "" {
		res = append(res, labelKeysFromDescription(p.Description, d)...)
	}

	return res
}

// LabelValues returns the values for a protocol
func (d *DefaultLabelStrategy) LabelValues(p *protocol.Protocol) []string {
	res := []string{p.Name, protoString(p), p.IPVersion, p.ImportFilter, p.ExportFilter}
	if d.descriptionLabels && p.Description != "" {
		res = append(res, labelValuesFromDescription(p.Description, d)...)
	}

	return res
}

func labelKeysFromDescription(desc string, d *DefaultLabelStrategy) (res []string) {
	reAllStringSubmatch := labelFindAllStringSubmatch(desc, d)
	for _, submatch := range reAllStringSubmatch {
		res = append(res, submatch[1])
	}

	return
}

func labelValuesFromDescription(desc string, d *DefaultLabelStrategy) (res []string) {
	reAllStringSubmatch := labelFindAllStringSubmatch(desc, d)
	for _, submatch := range reAllStringSubmatch {
		res = append(res, submatch[2])
	}

	return
}

func labelFindAllStringSubmatch(desc string, d *DefaultLabelStrategy) (result [][]string) {

	// Regex pattern captures "key: value" pair from the content.
	pattern := regexp.MustCompile(d.descriptionLabelsRegex)

	result = pattern.FindAllStringSubmatch(desc, -1)

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
	case protocol.RPKI:
		return "RPKI"
	case protocol.BFD:
		return "BFD"
	}

	return ""
}
