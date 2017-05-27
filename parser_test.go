package main

import (
	"testing"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/czerwonk/testutils/assert"
)

func TestEstablishedBgpOldTimeFormat(t *testing.T) {
	data := "foo    BGP      master   up     1481973060  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.Name, t)
	assert.IntEqual("proto", protocol.BGP, x.Proto, t)
	assert.IntEqual("established", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("filtered", 1, x.Filtered, t)
	assert.IntEqual("ipVersion", 4, x.IpVersion, t)
}

func TestEstablishedBgpCurrentTimeFormat(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.Name, t)
	assert.IntEqual("proto", protocol.BGP, x.Proto, t)
	assert.IntEqual("established", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("filtered", 1, x.Filtered, t)
	assert.IntEqual("ipVersion", 4, x.IpVersion, t)
	assert.IntEqual("uptime", 60, x.Uptime, t)
}

func TestIpv6Bgp(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 6)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.IntEqual("ipVersion", 6, x.IpVersion, t)
}

func TestActiveBgp(t *testing.T) {
	data := "bar    BGP      master   start   2016-01-01    Active\ntest\nbar"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "bar", x.Name, t)
	assert.IntEqual("proto", protocol.BGP, x.Proto, t)
	assert.IntEqual("established", 0, x.Up, t)
	assert.IntEqual("imported", 0, int(x.Imported), t)
	assert.IntEqual("exported", 0, int(x.Exported), t)
	assert.IntEqual("ipVersion", 4, x.IpVersion, t)
	assert.IntEqual("uptime", 0, int(x.Uptime), t)
}

func Test2BgpSessions(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nbar    BGP      master   start   2016-01-01    Active\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 2, len(p), t)
}

func TestOspfOldTimeFormat(t *testing.T) {
	data := "ospf1    OSPF      master   up     1481973060  Running\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "ospf1", x.Name, t)
	assert.IntEqual("proto", protocol.OSPF, x.Proto, t)
	assert.IntEqual("up", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.IntEqual("ipVersion", 4, x.IpVersion, t)
}

func TestOspfCurrentTimeFormat(t *testing.T) {
	data := "ospf1    OSPF      master   up     00:01:00  Running\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "ospf1", x.Name, t)
	assert.IntEqual("proto", protocol.OSPF, x.Proto, t)
	assert.IntEqual("up", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.IntEqual("ipVersion", 4, x.IpVersion, t)
	assert.IntEqual("uptime", 60, x.Uptime, t)
}

func TestOspfProtocolDown(t *testing.T) {
	data := "o_hrz    OSPF     t_hrz    down   1494926415  \n  Preference:     150\n  Input filter:   ACCEPT\n  Output filter:  REJECT\nxxx"
	p := parseOutput([]byte(data), 6)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "o_hrz", x.Name, t)
	assert.IntEqual("proto", protocol.OSPF, x.Proto, t)
	assert.IntEqual("up", 0, x.Up, t)
	assert.Int64Equal("imported", 0, x.Imported, t)
	assert.Int64Equal("exported", 0, x.Exported, t)
	assert.IntEqual("ipVersion", 6, x.IpVersion, t)
}

func TestOspfRunning(t *testing.T) {
	data := "ospf1    OSPF      master   up     00:01:00  Running\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.Float64Equal("running", 1, x.Attributes["running"], t)
}

func TestOspfAlone(t *testing.T) {
	data := "ospf1    OSPF      master   up     00:01:00  Alone\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.Float64Equal("running", 0, x.Attributes["running"], t)
}
