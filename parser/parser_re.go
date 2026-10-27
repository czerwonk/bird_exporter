package parser

import (
	"regexp"

	"github.com/czerwonk/bird_exporter/protocol"
)

var (
	protocolRegex    *regexp.Regexp
	descriptionRegex *regexp.Regexp
	messageRegex     *regexp.Regexp
	routeRegex       *regexp.Regexp
	routeChangeRegex *regexp.Regexp
	filterRegex      *regexp.Regexp
	channelRegex     *regexp.Regexp
)

var regexMatchers lineMatchers

func init() {
	protocolRegex = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+(MRT|BGP|BFD|OSPF|RPKI|RIP|RAdv|Pipe|Perf|Direct|Babel|Device|Kernel|Static)\s+([^\s]+)\s+([^\s]+)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:|\.\d+)|[^\s]+)(?:\s+(.*?))?$`)
	descriptionRegex = regexp.MustCompile(`^(?:1006-)?\s+Description:\s+(.*)`)
	messageRegex = regexp.MustCompile(`^(?:1006-)?\s+Message:`)
	routeRegex = regexp.MustCompile(`^\s+Routes:\s+(\d+) imported, (?:(\d+) filtered, )?(\d+) exported(?:, (\d+) preferred)?`)
	routeChangeRegex = regexp.MustCompile(`^(?:1006-)?\s+(Import|Export) (updates|withdraws):\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)(?:\s+(\d+|---)\s+(\d+|---))?`)
	filterRegex = regexp.MustCompile(`^(?:1006-)?\s+(Input|Output) filter:\s+(.*)`)
	channelRegex = regexp.MustCompile(`^\s+Channel (?:ipv|roa)(4|6)(?:\s|$)`)

	regexMatchers = lineMatchers{
		protocol:    protocolRegex.FindStringSubmatch,
		description: descriptionRegex.FindStringSubmatch,
		message:     messageRegex.FindStringSubmatch,
		channel:     channelRegex.FindStringSubmatch,
		routes:      routeRegex.FindStringSubmatch,
		routeChange: routeChangeRegex.FindStringSubmatch,
		filter:      filterRegex.FindStringSubmatch,
	}
}

// ParseProtocols parses bird output and returns protocol.Protocol structs
func ParseProtocols(data []byte, ipVersion string) []*protocol.Protocol {
	protocols, _ := ParseProtocolsWithError(data, ipVersion)
	return protocols
}

// ParseProtocolsWithError parses BIRD output and reports scanner failures.
func ParseProtocolsWithError(data []byte, ipVersion string) ([]*protocol.Protocol, error) {
	return parseProtocols(data, ipVersion, &regexMatchers)
}
