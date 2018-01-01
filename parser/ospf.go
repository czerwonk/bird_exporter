package parser

import (
	"regexp"

	"bufio"
	"bytes"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

type ospfRegex struct {
	area     *regexp.Regexp
	counters *regexp.Regexp
}

type ospfContext struct {
	line    string
	areas   []*protocol.OspfArea
	current *protocol.OspfArea
}

func init() {
	ospf = &ospfRegex{
		area:     regexp.MustCompile("Area: [^\\s]+ \\(([^\\s]+)\\)"),
		counters: regexp.MustCompile("Number of ([^:]+):\\s*(\\d+)"),
	}
}

var ospf *ospfRegex

func ParseOspf(data []byte) []*protocol.OspfArea {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &ospfContext{
		areas: make([]*protocol.OspfArea, 0),
	}

	for scanner.Scan() {
		c.line = strings.Trim(scanner.Text(), " ")
		parseLineForOspfArea(c)
		parseLineForOspfCounters(c)
	}

	return c.areas
}

func parseLineForOspfArea(c *ospfContext) {
	m := ospf.area.FindStringSubmatch(c.line)
	if m == nil {
		return
	}

	a := &protocol.OspfArea{Name: m[1]}
	c.current = a
	c.areas = append(c.areas, a)
}

func parseLineForOspfCounters(c *ospfContext) {
	if c.current == nil {
		return
	}

	m := ospf.counters.FindStringSubmatch(c.line)
	if m == nil {
		return
	}

	name := m[1]
	value := parseInt(m[2])

	if name == "interfaces" {
		c.current.InterfaceCount = value
	}

	if name == "neighbors" {
		c.current.NeighborCount = value
	}

	if name == "adjacent neighbors" {
		c.current.NeighborAdjacentCount = value
	}
}
