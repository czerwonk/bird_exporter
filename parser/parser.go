package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

const maxProtocolLineBytes = 1 << 20

var uptimeRegex = regexp.MustCompile(`^(?:((\d+):(\d{2}):(\d{2})(?:\.\d+)?)|(\d+)|(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)?))$`)

// Fixed words of route change and filter lines, shared by both parsers.
const (
	routeChangeImport    = "Import"
	routeChangeExport    = "Export"
	routeChangeUpdates   = "updates"
	routeChangeWithdraws = "withdraws"
	routeChangeUnset     = "---"
	filterInput          = "Input"
	filterOutput         = "Output"
	filterKey            = " filter:"
)

type context struct {
	current   *protocol.Protocol
	line      string
	handled   bool
	protocols []*protocol.Protocol
	ipVersion string
	match     *lineMatchers
}

// lineMatchers return submatches in the layout of the regular expressions in
// parser_re.go, or nil when the line does not match.
type lineMatchers struct {
	protocol    func(string) []string
	description func(string) []string
	message     func(string) []string
	channel     func(string) []string
	routes      func(string) []string
	routeChange func(string) []string
	filter      func(string) []string
}

func parseProtocols(data []byte, ipVersion string, match *lineMatchers) ([]*protocol.Protocol, error) {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64<<10), maxProtocolLineBytes)

	c := &context{protocols: make([]*protocol.Protocol, 0), ipVersion: ipVersion, match: match}

	var handlers = []func(*context){
		handleEmptyLine,
		parseLineForProtocol,
		skipLineForMessage,
		parseLineForDescription,
		parseLineForChannel,
		parseLineForRoutes,
		parseLineForRouteChanges,
		parseLineForFilterName,
	}

	for scanner.Scan() {
		c.line = strings.TrimRight(scanner.Text(), " ")
		c.handled = false

		for _, h := range handlers {
			if !c.handled {
				h(c)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return c.protocols, fmt.Errorf("scan BIRD protocols reply: %w", err)
	}

	return c.protocols, nil
}

func handleEmptyLine(c *context) {
	if c.line != "" {
		return
	}

	c.current = nil
	c.handled = true
}

func parseLineForProtocol(c *context) {
	match := c.match.protocol(c.line)

	if match == nil {
		return
	}
	proto := parseProto(match[2])
	ut := parseUptime(match[5])

	c.current = protocol.NewProtocol(match[1], proto, c.ipVersion, ut)
	c.current.Up = parseState(match[4])
	c.current.State = match[6]

	c.protocols = append(c.protocols, c.current)
	c.handled = true
}

// skipLineForMessage consumes the protocol message line. It can carry text
// received from a BGP peer (shutdown communication), so no other handler may
// interpret it.
func skipLineForMessage(c *context) {
	if c.match.message(c.line) == nil {
		return
	}

	c.handled = true
}

func parseLineForDescription(c *context) {
	if c.current == nil {
		return
	}

	match := c.match.description(c.line)

	if match == nil {
		return
	}

	if len(match) <= 1 {
		return
	}

	c.current.Description = strings.Join(match[1:], " ")
}

func parseProto(val string) protocol.Proto {
	switch val {
	case "BGP":
		return protocol.BGP
	case "OSPF":
		return protocol.OSPF
	case "Direct":
		return protocol.Direct
	case "Kernel":
		return protocol.Kernel
	case "Static":
		return protocol.Static
	case "Babel":
		return protocol.Babel
	case "RPKI":
		return protocol.RPKI
	case "BFD":
		return protocol.BFD
	}

	return protocol.PROTO_UNKNOWN
}

func parseState(state string) int {
	if state == "up" {
		return 1
	}

	return 0
}

func parseUptime(value string) int {
	match := uptimeRegex.FindStringSubmatch(value)

	if match == nil {
		return 0
	}

	if len(match[1]) > 0 {
		return parseUptimeForDuration(match)
	}

	if len(match[5]) > 0 {
		return parseUptimeForTimestamp(value)
	}

	return parseUptimeForIso(value)
}

func parseLineForChannel(c *context) {
	if c.ipVersion != "" || c.current == nil {
		return
	}

	channel := c.match.channel(c.line)
	if channel == nil {
		return
	}

	if len(c.current.IPVersion) == 0 {
		c.current.IPVersion = channel[1]
	} else {
		c.current = &protocol.Protocol{
			Name:        c.current.Name,
			Description: c.current.Description,
			Proto:       c.current.Proto,
			Up:          c.current.Up,
			State:       c.current.State,
			Uptime:      c.current.Uptime,
			IPVersion:   channel[1],
		}
		c.protocols = append(c.protocols, c.current)
	}

	c.handled = true
}

func parseLineForRoutes(c *context) {
	if c.current == nil {
		return
	}

	match := c.match.routes(c.line)
	if match == nil {
		return
	}

	c.current.Imported, _ = strconv.ParseInt(match[1], 10, 64)
	c.current.Exported, _ = strconv.ParseInt(match[3], 10, 64)

	if len(match[2]) > 0 {
		c.current.Filtered, _ = strconv.ParseInt(match[2], 10, 64)
	}

	if len(match[4]) > 0 {
		c.current.Preferred, _ = strconv.ParseInt(match[4], 10, 64)
	}

	c.handled = true
}

func parseLineForRouteChanges(c *context) {
	if c.current == nil {
		return
	}

	match := c.match.routeChange(c.line)
	if match == nil {
		return
	}
	c.handled = true

	x := getRouteChangeCount(match, c.current)
	x.Received = parseRouteChangeValue(match[3])
	x.Rejected = parseRouteChangeValue(match[4])
	x.Filtered = parseRouteChangeValue(match[5])
	x.Ignored = parseRouteChangeValue(match[6])

	if len(match[8]) > 0 { // v3 format
		x.RxLimit = parseRouteChangeValue(match[7])
		x.Limit = parseRouteChangeValue(match[8])
		x.Accepted = parseRouteChangeValue(match[9])
		c.current.RouteChangeFormatV3 = true
	} else {
		x.Accepted = parseRouteChangeValue(match[7])
	}
}

func getRouteChangeCount(values []string, p *protocol.Protocol) *protocol.RouteChangeCount {
	if values[1] == routeChangeImport {
		if values[2] == routeChangeUpdates {
			return &p.ImportUpdates
		}

		return &p.ImportWithdraws
	}

	if values[2] == routeChangeUpdates {
		return &p.ExportUpdates
	}

	return &p.ExportWithdraws
}

func parseRouteChangeValue(value string) int64 {
	if value == routeChangeUnset {
		return 0
	}

	return parseInt(value)
}

func parseLineForFilterName(c *context) {
	if c.current == nil {
		return
	}

	match := c.match.filter(c.line)
	if match == nil {
		return
	}

	if match[1] == filterInput {
		c.current.ImportFilter = match[2]
	} else {
		c.current.ExportFilter = match[2]
	}

	c.handled = true
}
