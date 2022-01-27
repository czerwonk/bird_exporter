package parser

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

var (
	bfdProtocolNameRegex *regexp.Regexp
	bfdSessionRegex      *regexp.Regexp
)

func init() {
	bfdProtocolNameRegex = regexp.MustCompile(`^([^\s]+):$`)
	bfdSessionRegex = regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+(Up|Down)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)\s+([0-9\.]+)\s+([0-9\.]+)$`)
}

type bfdContext struct {
	line     string
	sessions []*protocol.BFDSession
	protocol string
	handled  bool
}

func ParseBFDSessions(data []byte) []*protocol.BFDSession {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &bfdContext{
		sessions: make([]*protocol.BFDSession, 0),
	}

	for scanner.Scan() {
		c.line = strings.TrimSpace(scanner.Text())
		c.handled = false

		parseBFDProtocolLine(c)
		parseBFDSessionLine(c)
	}

	return c.sessions
}

func parseBFDProtocolLine(c *bfdContext) {
	if c.handled {
		return
	}

	m := bfdProtocolNameRegex.FindStringSubmatch(c.line)
	if m == nil {
		return
	}

	c.protocol = m[1]
	c.handled = true
}

func parseBFDSessionLine(c *bfdContext) {
	if c.handled {
		return
	}

	m := bfdSessionRegex.FindStringSubmatch(c.line)
	if m == nil {
		return
	}

	sess := protocol.BFDSession{
		ProtocolName: c.protocol,
		IP:           m[1],
		Interface:    m[2],
		Since:        parseUptime(m[4]),
		Interval:     parseFloat(m[5]),
		Timeout:      parseFloat(m[6]),
	}

	if m[3] == "Up" {
		sess.Up = true
	}

	c.sessions = append(c.sessions, &sess)
	c.handled = true
}
