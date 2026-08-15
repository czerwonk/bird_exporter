package metrics

import (
	"errors"
	"regexp"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

var prometheusLabelNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

const (
	maxDescriptionLabels          = 32
	maxDescriptionLabelNameBytes  = 128
	maxDescriptionLabelValueBytes = 1024
	prometheusReservedLabelPrefix = "__"
)

var baseProtocolLabelNames = map[string]struct{}{
	"name":          {},
	"proto":         {},
	"ip_version":    {},
	"import_filter": {},
	"export_filter": {},
	"state":         {},
}

// DefaultLabelStrategy defines the labels to add to a metric and its data retrieval method.
type DefaultLabelStrategy struct {
	descriptionLabels      bool
	descriptionLabelsRegex *regexp.Regexp
}

// ValidateDescriptionLabelsRegex validates startup configuration without panicking.
func ValidateDescriptionLabelsRegex(descriptionLabels bool, expression string) error {
	if !descriptionLabels {
		return nil
	}

	compiled, err := regexp.Compile(expression)
	if err != nil {
		return err
	}
	if compiled.NumSubexp() != 2 {
		return errors.New("description labels regex must contain exactly two capture groups")
	}

	return nil
}

// NewDefaultLabelStrategy creates a label strategy. Invalid optional
// configuration disables description labels instead of panicking.
func NewDefaultLabelStrategy(descriptionLabels bool, expression string) *DefaultLabelStrategy {
	if err := ValidateDescriptionLabelsRegex(descriptionLabels, expression); err != nil {
		return &DefaultLabelStrategy{}
	}
	if !descriptionLabels {
		return &DefaultLabelStrategy{}
	}

	return &DefaultLabelStrategy{
		descriptionLabels:      true,
		descriptionLabelsRegex: regexp.MustCompile(expression),
	}
}

// LabelNames returns the list of label names.
func (d *DefaultLabelStrategy) LabelNames(p *protocol.Protocol) []string {
	result := []string{"name", "proto", "ip_version", "import_filter", "export_filter"}
	names, _ := d.labelsFromDescription(p)
	return append(result, names...)
}

// LabelValues returns values in the same order as LabelNames.
func (d *DefaultLabelStrategy) LabelValues(p *protocol.Protocol) []string {
	result := []string{p.Name, protoString(p), p.IPVersion, p.ImportFilter, p.ExportFilter}
	_, values := d.labelsFromDescription(p)
	return append(result, values...)
}

func (d *DefaultLabelStrategy) labelsFromDescription(p *protocol.Protocol) ([]string, []string) {
	if !d.descriptionLabels || d.descriptionLabelsRegex == nil || p.Description == "" {
		return nil, nil
	}

	matches := d.descriptionLabelsRegex.FindAllStringSubmatch(p.Description, maxDescriptionLabels+1)
	if len(matches) > maxDescriptionLabels {
		return nil, nil
	}
	names := make([]string, 0, len(matches))
	values := make([]string, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))

	for _, submatch := range matches {
		if len(submatch) != 3 {
			return nil, nil
		}

		name := strings.TrimSpace(submatch[1])
		value := strings.TrimSpace(submatch[2])
		if len(name) > maxDescriptionLabelNameBytes || len(value) > maxDescriptionLabelValueBytes {
			return nil, nil
		}
		if !prometheusLabelNameRegex.MatchString(name) || strings.HasPrefix(name, prometheusReservedLabelPrefix) {
			return nil, nil
		}
		if _, reserved := baseProtocolLabelNames[name]; reserved {
			return nil, nil
		}
		if _, duplicate := seen[name]; duplicate {
			return nil, nil
		}

		seen[name] = struct{}{}
		names = append(names, name)
		values = append(values, value)
	}

	return names, values
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
