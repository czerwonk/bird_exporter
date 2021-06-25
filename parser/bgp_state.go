package parser

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

var (
	nameRegex     *regexp.Regexp
	bgpStateRegex *regexp.Regexp
)

type bgpStateContext struct {
	line    string
	current *protocol.BgpState
}

func init() {
	bgpStateRegex = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+(MRT|BGP|BFD|OSPF|RPKI|RIP|RAdv|Pipe|Perf|Direct|Babel|Device|Kernel|Static)\s+([^\s]+)\s+([^\s]+)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)\s+(Idle|Connect|Active|OpenSent|OpenConfirm|Established|Close)(?:\s+(.*?))?$`)
}

func ParseBgpState(data []byte) *protocol.BgpState {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &bgpStateContext{
		current: nil,
	}

	for scanner.Scan() {
		c.line = strings.TrimRight(scanner.Text(), " ")
		if c.line == "" {
			c.current = nil
		}
		m := bgpStateRegex.FindStringSubmatch(c.line)
		if m != nil {
			s := &protocol.BgpState{Name: m[1]}
			c.current = s
			c.current.State = m[6]
			break
		}
	}

	return c.current
}
