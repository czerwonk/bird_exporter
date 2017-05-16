package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

var (
	protocolRegex *regexp.Regexp
	routeRegex    *regexp.Regexp
	uptimeRegex   *regexp.Regexp
)

func init() {
	protocolRegex, _ = regexp.Compile("^([^\\s]+)\\s+(BGP|OSPF)\\s+([^\\s]+)\\s+([^\\s]+)\\s+([^\\s]+)\\s+(.*?)\\s*$")
	routeRegex, _ = regexp.Compile("^\\s+Routes:\\s+(\\d+) imported, (?:\\d+ filtered, )?(\\d+) exported")
	uptimeRegex, _ = regexp.Compile("^(?:((\\d+):(\\d{2}):(\\d{2}))|\\d+)$")
}

func parseOutput(data []byte, ipVersion int) []*protocol {
	protocols := make([]*protocol, 0)

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	var current *protocol = nil

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

func parseLineForProtocol(line string, ipVersion int) (*protocol, bool) {
	match := protocolRegex.FindStringSubmatch(line)

	if match == nil {
		return nil, false
	}

	proto := parseProto(match[2])
	up := parseState(match[4], proto)
	ut := parseUptime(match[5])
	p := &protocol{proto: proto, name: match[1], ipVersion: ipVersion, up: up, uptime: ut, attributes: make(map[string]interface{})}

	return p, true
}

func parseProto(val string) int {
	switch val {
	case "BGP":
		return BGP
	case "OSPF":
		return OSPF
	}

	return PROTO_UNKNOWN
}

func parseLineForRoutes(line string, p *protocol) {
	match := routeRegex.FindStringSubmatch(line)

	if match != nil {
		p.imported, _ = strconv.ParseInt(match[1], 10, 64)
		p.exported, _ = strconv.ParseInt(match[2], 10, 64)
	}
}

func parseState(state string, proto int) int {
	if state == "up" {
		return 1
	} else {
		return 0
	}
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
		log.Println(err)
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
		log.Println(err)
		return 0
	}

	return i
}
