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
192.168.64.9              enp0s2     Up         2022-01-27 09:00:00 1697620076    0.100    1.000
192.168.64.10             enp0s2     Down       2022-01-27 08:00:00    0.300    0.000
192.168.64.12             enp0s2     Init       2022-01-27 08:00:00    0.300    5.000`

	s := ParseBFDSessions("bfd1", []byte(data))

	assert.Equal(t, 3, len(s), "session count")

	s1 := protocol.BFDSession{
		ProtocolName: "bfd1",
		IP:           "192.168.64.9",
		Interface:    "enp0s2",
		Up:           true,
		Since:        3600,
		SinceEpoch:   1697620076,
		Interval:     0.1,
		Timeout:      1,
	}
	s2 := protocol.BFDSession{
		ProtocolName: "bfd1",
		IP:           "192.168.64.10",
		Interface:    "enp0s2",
		Up:           false,
		Since:        7200,
		SinceEpoch:   0,
		Interval:     0.3,
		Timeout:      0,
	}
	s3 := protocol.BFDSession{
		ProtocolName: "bfd1",
		IP:           "192.168.64.12",
		Interface:    "enp0s2",
		Up:           false,
		Since:        7200,
		SinceEpoch:   0,
		Interval:     0.3,
		Timeout:      5,
	}
	assert.Equal(t, []*protocol.BFDSession{&s1, &s2, &s3}, s, "sessions")
}
