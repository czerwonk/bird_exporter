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
	assert.Int64Equal("preferred", 100, x.Preferred, t)
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
	assert.Int64Equal("preferred", 100, x.Preferred, t)
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

func TestUpdateAndWithdrawCounts(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\n" +
		"  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\n" +
		"  Route change stats:     received   rejected   filtered    ignored   accepted\n" +
		"    Import updates:              1          2          3          4          5\n" +
		"    Import withdraws:            6          7          8          9         10\n" +
		"    Export updates:             11         12         13         14         15\n" +
		"    Export withdraws:           16         17         18         19        ---"
	p := parseOutput([]byte(data), 4)
	x := p[0]

	assert.Int64Equal("import updates received", 1, x.ImportUpdates.Received, t)
	assert.Int64Equal("import updates rejected", 2, x.ImportUpdates.Rejected, t)
	assert.Int64Equal("import updates filtered", 3, x.ImportUpdates.Filtered, t)
	assert.Int64Equal("import updates ignored", 4, x.ImportUpdates.Ignored, t)
	assert.Int64Equal("import updates accepted", 5, x.ImportUpdates.Accepted, t)
	assert.Int64Equal("import withdraws received", 6, x.ImportWithdraws.Received, t)
	assert.Int64Equal("import withdraws rejected", 7, x.ImportWithdraws.Rejected, t)
	assert.Int64Equal("import withdraws filtered", 8, x.ImportWithdraws.Filtered, t)
	assert.Int64Equal("import withdraws ignored", 9, x.ImportWithdraws.Ignored, t)
	assert.Int64Equal("import withdraws accepted", 10, x.ImportWithdraws.Accepted, t)
	assert.Int64Equal("export updates received", 11, x.ExportUpdates.Received, t)
	assert.Int64Equal("export updates rejected", 12, x.ExportUpdates.Rejected, t)
	assert.Int64Equal("export updates filtered", 13, x.ExportUpdates.Filtered, t)
	assert.Int64Equal("export updates ignored", 14, x.ExportUpdates.Ignored, t)
	assert.Int64Equal("export updates accepted", 15, x.ExportUpdates.Accepted, t)
	assert.Int64Equal("export withdraws received", 16, x.ExportWithdraws.Received, t)
	assert.Int64Equal("export withdraws rejected", 17, x.ExportWithdraws.Rejected, t)
	assert.Int64Equal("export withdraws filtered", 18, x.ExportWithdraws.Filtered, t)
	assert.Int64Equal("export withdraws ignored", 19, x.ExportWithdraws.Ignored, t)
	assert.Int64Equal("export withdraws accepted", 0, x.ExportWithdraws.Accepted, t)
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
	assert.Int64Equal("preferred", 100, x.Preferred, t)
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
	assert.Int64Equal("preferred", 100, x.Preferred, t)
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
