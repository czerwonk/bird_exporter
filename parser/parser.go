package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/common/log"
)

var (
	protocolRegex    *regexp.Regexp
	routeRegex       *regexp.Regexp
	uptimeRegex      *regexp.Regexp
	routeChangeRegex *regexp.Regexp
	channelRegex     *regexp.Regexp
)

type context struct {
	current   *protocol.Protocol
	line      string
	handled   bool
	protocols []*protocol.Protocol
	ipVersion string
}

func init() {
	protocolRegex = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+(BGP|OSPF|Direct|Device|Kernel|Static)\s+([^\s]+)\s+([^\s]+)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)(?:\s+(.*?))?$`)
	routeRegex = regexp.MustCompile(`^\s+Routes:\s+(\d+) imported, (?:(\d+) filtered, )?(\d+) exported(?:, (\d+) preferred)?`)
	uptimeRegex = regexp.MustCompile(`^(?:((\d+):(\d{2}):(\d{2}))|(\d+)|(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}))$`)
	routeChangeRegex = regexp.MustCompile(`(Import|Export) (updates|withdraws):\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s*`)
	channelRegex = regexp.MustCompile(`Channel ipv(4|6)`)
}

// Parser parses bird output and returns protocol.Protocol structs
func ParseProtocols(data []byte, ipVersion string) []*protocol.Protocol {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &context{protocols: make([]*protocol.Protocol, 0), ipVersion: ipVersion}

	var handlers = []func(*context){
		handleEmptyLine,
		parseLineForProtocol,
		parseLineForChannel,
		parseLineForRoutes,
		parseLineForRouteChanges,
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

	return c.protocols
}

func handleEmptyLine(c *context) {
	if c.line != "" {
		return
	}

	c.current = nil
	c.handled = true
}

func parseLineForProtocol(c *context) {
	match := protocolRegex.FindStringSubmatch(c.line)

	if match == nil {
		return
	}

	proto := parseProto(match[2])
	ut := parseUptime(match[5])

	c.current = protocol.NewProtocol(match[1], proto, c.ipVersion, ut)
	c.current.Up = parseState(match[4])

	fillAttributes(c.current, match)

	c.protocols = append(c.protocols, c.current)
	c.handled = true
}

func parseProto(val string) int {
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

func parseUptimeForIso(s string) int {
	start, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		log.Errorln(err)
		return 0
	}

	return int(time.Since(start).Seconds())
}

func parseUptimeForDuration(duration []string) int {
	h := parseInt(duration[2])
	m := parseInt(duration[3])
	s := parseInt(duration[4])
	str := fmt.Sprintf("%dh%dm%ds", h, m, s)

	d, err := time.ParseDuration(str)
	if err != nil {
		log.Errorln(err)
		return 0
	}

	return int(d.Seconds())
}

func parseUptimeForTimestamp(timestamp string) int {
	since := parseInt(timestamp)

	s := time.Unix(since, 0)
	d := time.Since(s)
	return int(d.Seconds())
}

func parseLineForChannel(c *context) {
	if c.ipVersion != "" || c.current == nil {
		return
	}

	channel := channelRegex.FindStringSubmatch(c.line)
	if channel == nil {
		return
	}

	if len(c.current.IpVersion) == 0 {
		c.current.IpVersion = channel[1]
	} else {
		c.current = &protocol.Protocol{
			Name:      c.current.Name,
			Proto:     c.current.Proto,
			Up:        c.current.Up,
			Uptime:    c.current.Uptime,
			IpVersion: channel[1],
		}
		c.protocols = append(c.protocols, c.current)
	}

	c.handled = true
}

func parseLineForRoutes(c *context) {
	if c.current == nil {
		return
	}

	match := routeRegex.FindStringSubmatch(c.line)
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

	match := routeChangeRegex.FindStringSubmatch(c.line)
	if match == nil {
		return
	}

	x := getRouteChangeCount(match, c.current)
	x.Received = parseRouteChangeValue(match[3])
	x.Rejected = parseRouteChangeValue(match[4])
	x.Filtered = parseRouteChangeValue(match[5])
	x.Ignored = parseRouteChangeValue(match[6])
	x.Accepted = parseRouteChangeValue(match[7])

	c.handled = true
}

func getRouteChangeCount(values []string, p *protocol.Protocol) *protocol.RouteChangeCount {
	if values[1] == "Import" {
		if values[2] == "updates" {
			return &p.ImportUpdates
		}

		return &p.ImportWithdraws
	} else {
		if values[2] == "updates" {
			return &p.ExportUpdates
		}

		return &p.ExportWithdraws
	}
}

func parseRouteChangeValue(value string) int64 {
	if value == "---" {
		return 0
	}

	return parseInt(value)
}

func parseInt(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		log.Errorln(err)
		return 0
	}

	return i
}

func fillAttributes(p *protocol.Protocol, m []string) {
	if p.Proto == protocol.OSPF {
		p.Attributes["running"] = float64(parseOspfRunning(m[6]))
	}
}

func parseOspfRunning(state string) int {
	if state == "Running" {
		return 1
	}

	return 0
}
