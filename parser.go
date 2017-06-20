package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/common/log"
)

var (
	protocolRegex *regexp.Regexp
	routeRegex    *regexp.Regexp
	uptimeRegex   *regexp.Regexp
)

func init() {
	protocolRegex, _ = regexp.Compile("^(?:1002\\-)?([^\\s]+)\\s+(BGP|OSPF)\\s+([^\\s]+)\\s+([^\\s]+)\\s+([^\\s]+)\\s+(.*?)\\s*$")
	routeRegex, _ = regexp.Compile("^\\s+Routes:\\s+(\\d+) imported, (?:(\\d+) filtered, )?(\\d+) exported")
	uptimeRegex, _ = regexp.Compile("^(?:((\\d+):(\\d{2}):(\\d{2}))|\\d+)$")
}

func parseOutput(data []byte, ipVersion int) []*protocol.Protocol {
	protocols := make([]*protocol.Protocol, 0)

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	var current *protocol.Protocol = nil

	for scanner.Scan() {
		line := scanner.Text()
		if p, ok := parseLineForProtocol(line, ipVersion); ok {
			current = p
			protocols = append(protocols, current)
		}

		if current != nil {
			parseLineForRoutes(line, current)
		}

		if line == "" {
			current = nil
		}
	}

	return protocols
}

func parseLineForProtocol(line string, ipVersion int) (*protocol.Protocol, bool) {
	match := protocolRegex.FindStringSubmatch(line)

	if match == nil {
		return nil, false
	}

	proto := parseProto(match[2])
	up := parseState(match[4])
	ut := parseUptime(match[5])
	p := &protocol.Protocol{Proto: proto, Name: match[1], IpVersion: ipVersion, Up: up, Uptime: ut, Attributes: make(map[string]float64)}
	fillAttributes(p, match)

	return p, true
}

func parseProto(val string) int {
	switch val {
	case "BGP":
		return protocol.BGP
	case "OSPF":
		return protocol.OSPF
	}

	return protocol.PROTO_UNKNOWN
}

func parseLineForRoutes(line string, p *protocol.Protocol) {
	match := routeRegex.FindStringSubmatch(line)

	if match != nil {
		p.Imported, _ = strconv.ParseInt(match[1], 10, 64)
		p.Exported, _ = strconv.ParseInt(match[3], 10, 64)

		if len(match[2]) > 0 {
			p.Filtered, _ = strconv.ParseInt(match[2], 10, 64)
		}
	}
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

	if match[1] != "" {
		return parseUptimeForDuration(match)
	}

	return parseUptimeForTimestamp(value)
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
