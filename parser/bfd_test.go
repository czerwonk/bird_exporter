package parser

import (
	"testing"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/stretchr/testify/assert"
)

func TestParseBFDSessions(t *testing.T) {
	overrideNowFunc(func() time.Time {
		return time.Date(2022, 1, 27, 10, 0, 0, 0, time.Local)
	})

	data := `BIRD 2.0.7 ready.
bfd1:
IP address                Interface  State      Since         Interval  Timeout
192.168.64.9              enp0s2     Up         2022-01-27 09:00:00    0.100    1.000
192.168.64.10             enp0s2     Down       2022-01-27 08:00:00    0.300    0.000`

	s := ParseBFDSessions([]byte(data))

	assert.Equal(t, 2, len(s), "session count")

	s1 := protocol.BFDSession{
		ProtocolName: "bfd1",
		IP:           "192.168.64.9",
		Interface:    "enp0s2",
		Up:           true,
		Since:        3600,
		Interval:     0.1,
		Timeout:      1,
	}
	s2 := protocol.BFDSession{
		ProtocolName: "bfd1",
		IP:           "192.168.64.10",
		Interface:    "enp0s2",
		Up:           false,
		Since:        7200,
		Interval:     0.3,
		Timeout:      0,
	}
	assert.Equal(t, []*protocol.BFDSession{&s1, &s2}, s, "sessions")
}
