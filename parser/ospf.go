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
	areas   []*protocol.OSPFArea
	current *protocol.OSPFArea
}

func init() {
	ospf = &ospfRegex{
		area:     regexp.MustCompile("Area: [^\\s]+ \\(([^\\s]+)\\)"),
		counters: regexp.MustCompile("Number of ([^:]+):\\s*(\\d+)"),
	}
}

var ospf *ospfRegex

func ParseOSPF(data []byte) []*protocol.OSPFArea {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &ospfContext{
		areas: make([]*protocol.OSPFArea, 0),
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

	a := &protocol.OSPFArea{Name: m[1]}
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
