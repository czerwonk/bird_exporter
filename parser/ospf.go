package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
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
		area:     regexp.MustCompile(`Area: [^\s]+ \(([^\s]+)\)`),
		counters: regexp.MustCompile(`Number of ([^:]+):\s*(\d+)`),
	}
}

var ospf *ospfRegex

func ParseOSPF(data []byte) []*protocol.OSPFArea {
	areas, _ := ParseOSPFWithError(data)
	return areas
}

// ParseOSPFWithError parses OSPF area output and reports scanner failures.
func ParseOSPFWithError(data []byte) ([]*protocol.OSPFArea, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64<<10), maxProtocolLineBytes)

	c := &ospfContext{
		areas: make([]*protocol.OSPFArea, 0),
	}

	for scanner.Scan() {
		c.line = strings.Trim(scanner.Text(), " ")
		parseLineForOspfArea(c)
		parseLineForOspfCounters(c)
	}
	if err := scanner.Err(); err != nil {
		return c.areas, fmt.Errorf("scan BIRD OSPF reply: %w", err)
	}

	return c.areas, nil
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
