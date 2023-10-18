package parser

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

var (
	bfdSessionRegex *regexp.Regexp
)

func init() {
	bfdSessionRegex = regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+(Up|Down|Init)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)\s+(\d{1,})?\s+([0-9\.]+)\s+([0-9\.]+)$`)
}

type bfdContext struct {
	line     string
	sessions []*protocol.BFDSession
	protocol string
}

func ParseBFDSessions(protocolName string, data []byte) []*protocol.BFDSession {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	c := &bfdContext{
		sessions: make([]*protocol.BFDSession, 0),
		protocol: protocolName,
	}

	for scanner.Scan() {
		c.line = strings.TrimSpace(scanner.Text())
		parseBFDSessionLine(c)
	}

	return c.sessions
}

func parseBFDSessionLine(c *bfdContext) {
	m := bfdSessionRegex.FindStringSubmatch(c.line)
	if m == nil {
		return
	}
	var since_epoch int64
	if m[5] == "" {
		since_epoch = 0
	} else {
		since_epoch = parseInt(m[5])
	}

	sess := protocol.BFDSession{
		ProtocolName: c.protocol,
		IP:           m[1],
		Interface:    m[2],
		Since:        parseUptime(m[4]),
		SinceEpoch:   since_epoch,
		Interval:     parseFloat(m[6]),
		Timeout:      parseFloat(m[7]),
	}

	if m[3] == "Up" {
		sess.Up = true
	}

	c.sessions = append(c.sessions, &sess)
}
