package parser

import (
	"testing"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/czerwonk/testutils/assert"
)

func TestEstablishedBgpOldTimeFormat(t *testing.T) {
	overrideNowFunc(func() time.Time {
		return time.Date(2018, 1, 1, 2, 0, 0, 0, time.UTC)
	})

	data := "foo    BGP      master   up     1514768400  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "4")

	x := p[0]
	assert.StringEqual("name", "foo", x.Name, t)
	assert.IntEqual("proto", int(protocol.BGP), int(x.Proto), t)
	assert.IntEqual("established", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("filtered", 1, x.Filtered, t)
	assert.Int64Equal("preferred", 100, x.Preferred, t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
	assert.Int64Equal("uptime", 3600, int64(x.Uptime), t)
}

func TestEstablishedBgpCurrentTimeFormat(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "4")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.Name, t)
	assert.IntEqual("proto", int(protocol.BGP), int(x.Proto), t)
	assert.IntEqual("established", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("filtered", 1, x.Filtered, t)
	assert.Int64Equal("preferred", 100, x.Preferred, t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
	assert.IntEqual("uptime", 60, x.Uptime, t)
}

func TestEstablishedBgpIsoLongTimeFormat(t *testing.T) {
	overrideNowFunc(func() time.Time {
		return time.Date(2018, 1, 1, 2, 0, 0, 0, time.Local)
	})

	data := "foo    BGP      master   up     2018-01-01 01:00:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "4")

	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.Name, t)
	assert.IntEqual("proto", int(protocol.BGP), int(x.Proto), t)
	assert.IntEqual("established", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("filtered", 1, x.Filtered, t)
	assert.Int64Equal("preferred", 100, x.Preferred, t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
	assert.Int64Equal("uptime", 3600, int64(x.Uptime), t)
}

func TestIpv6BGP(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "6")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("ipVersion", "6", x.IPVersion, t)
}

func TestActiveBGP(t *testing.T) {
	data := "bar    BGP      master   start   2016-01-01    Active\ntest\nbar"
	p := ParseProtocols([]byte(data), "4")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "bar", x.Name, t)
	assert.IntEqual("proto", int(protocol.BGP), int(x.Proto), t)
	assert.IntEqual("established", 0, x.Up, t)
	assert.IntEqual("imported", 0, int(x.Imported), t)
	assert.IntEqual("exported", 0, int(x.Exported), t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
	assert.IntEqual("uptime", 0, int(x.Uptime), t)
}

func Test2BGPSessions(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nbar    BGP      master   start   2016-01-01    Active\nxxx"
	p := ParseProtocols([]byte(data), "4")
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
	p := ParseProtocols([]byte(data), "4")
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

func TestWithBird2(t *testing.T) {
	data := "Name       Proto      Table      State  Since         Info\n" +
		"bgp1       BGP        master     up     1494926415\n" +
		"  Channel ipv6\n" +
		"    Routes:         1 imported, 2 filtered, 3 exported, 4 preferred\n" +
		"    Input filter:   none\n" +
		"    Output filter:  all\n" +
		"\n" +
		"direct1    Direct     ---        up     1513027903\n" +
		"  Channel ipv4\n" +
		"    State:          UP\n" +
		"    Table:          master4\n" +
		"    Preference:     240\n" +
		"    Input filter:   ACCEPT\n" +
		"    Output filter:  REJECT\n" +
		"    Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\n" +
		"    Route change stats:     received   rejected   filtered    ignored   accepted\n" +
		"      Import updates:              1          2          3          4          5\n" +
		"      Import withdraws:            6          7          8          9         10\n" +
		"      Export updates:             11         12         13         14         15\n" +
		"      Export withdraws:           16         17         18         19        ---\n" +
		"  Channel ipv6\n" +
		"    State:          UP\n" +
		"    Table:          master6\n" +
		"    Preference:     240\n" +
		"    Input filter:   ACCEPT\n" +
		"    Output filter:  REJECT\n" +
		"    Routes:         3 imported, 7 filtered, 5 exported, 13 preferred\n" +
		"    Route change stats:     received   rejected   filtered    ignored   accepted\n" +
		"      Import updates:             20         21         22         23         24\n" +
		"      Import withdraws:           25         26         27         28         29\n" +
		"      Export updates:             30         31         32         33         34\n" +
		"      Export withdraws:           35         36         37         38        ---\n" +
		"\n" +
		"ospf1      OSPF       master     up     1494926415\n" +
		"  Channel ipv4\n" +
		"    Routes:         4 imported, 3 filtered, 2 exported, 1 preferred\n" +
		"\n"

	p := ParseProtocols([]byte(data), "")
	assert.IntEqual("protocols", 4, len(p), t)

	x := p[0]
	assert.StringEqual("BGP ipv6 name", "bgp1", x.Name, t)
	assert.IntEqual("BGP ipv6 proto", int(protocol.BGP), int(x.Proto), t)
	assert.StringEqual("BGP ipv6 ip version", "6", x.IPVersion, t)
	assert.Int64Equal("BGP ipv6 imported", 1, x.Imported, t)
	assert.Int64Equal("BGP ipv6 exported", 3, x.Exported, t)
	assert.Int64Equal("BGP ipv6 filtered", 2, x.Filtered, t)
	assert.Int64Equal("BGP ipv6 preferred", 4, x.Preferred, t)
	assert.StringEqual("BGP import filter", "none", x.ImportFilter, t)
	assert.StringEqual("BGP export filter", "all", x.ExportFilter, t)

	x = p[1]
	assert.StringEqual("Direct ipv4 name", "direct1", x.Name, t)
	assert.IntEqual("Direct ipv4 proto", int(protocol.Direct), int(x.Proto), t)
	assert.StringEqual("Direct ipv4 ip version", "4", x.IPVersion, t)
	assert.Int64Equal("Direct ipv4 imported", 12, x.Imported, t)
	assert.Int64Equal("Direct ipv4 exported", 34, x.Exported, t)
	assert.Int64Equal("Direct ipv4 filtered", 1, x.Filtered, t)
	assert.Int64Equal("Direct ipv4 preferred", 100, x.Preferred, t)
	assert.Int64Equal("Direct ipv4 import updates received", 1, x.ImportUpdates.Received, t)
	assert.Int64Equal("Direct ipv4 import updates rejected", 2, x.ImportUpdates.Rejected, t)
	assert.Int64Equal("Direct ipv4 import updates filtered", 3, x.ImportUpdates.Filtered, t)
	assert.Int64Equal("Direct ipv4 import updates ignored", 4, x.ImportUpdates.Ignored, t)
	assert.Int64Equal("Direct ipv4 import updates accepted", 5, x.ImportUpdates.Accepted, t)
	assert.Int64Equal("Direct ipv4 import withdraws received", 6, x.ImportWithdraws.Received, t)
	assert.Int64Equal("Direct ipv4 import withdraws rejected", 7, x.ImportWithdraws.Rejected, t)
	assert.Int64Equal("Direct ipv4 import withdraws filtered", 8, x.ImportWithdraws.Filtered, t)
	assert.Int64Equal("Direct ipv4 import withdraws ignored", 9, x.ImportWithdraws.Ignored, t)
	assert.Int64Equal("Direct ipv4 import withdraws accepted", 10, x.ImportWithdraws.Accepted, t)
	assert.Int64Equal("Direct ipv4 export updates received", 11, x.ExportUpdates.Received, t)
	assert.Int64Equal("Direct ipv4 export updates rejected", 12, x.ExportUpdates.Rejected, t)
	assert.Int64Equal("Direct ipv4 export updates filtered", 13, x.ExportUpdates.Filtered, t)
	assert.Int64Equal("Direct ipv4 export updates ignored", 14, x.ExportUpdates.Ignored, t)
	assert.Int64Equal("Direct ipv4 export updates accepted", 15, x.ExportUpdates.Accepted, t)
	assert.Int64Equal("Direct ipv4 export withdraws received", 16, x.ExportWithdraws.Received, t)
	assert.Int64Equal("Direct ipv4 export withdraws rejected", 17, x.ExportWithdraws.Rejected, t)
	assert.Int64Equal("Direct ipv4 export withdraws filtered", 18, x.ExportWithdraws.Filtered, t)
	assert.Int64Equal("Direct ipv4 export withdraws ignored", 19, x.ExportWithdraws.Ignored, t)
	assert.Int64Equal("Direct ipv4 export withdraws accepted", 0, x.ExportWithdraws.Accepted, t)

	x = p[2]
	assert.StringEqual("Direct ipv6 name", "direct1", x.Name, t)
	assert.IntEqual("Direct ipv6 proto", int(protocol.Direct), int(x.Proto), t)
	assert.StringEqual("Direct ipv6 ip version", "6", x.IPVersion, t)
	assert.Int64Equal("Direct ipv6 imported", 3, x.Imported, t)
	assert.Int64Equal("Direct ipv6 exported", 5, x.Exported, t)
	assert.Int64Equal("Direct ipv6 filtered", 7, x.Filtered, t)
	assert.Int64Equal("Direct ipv6 preferred", 13, x.Preferred, t)
	assert.Int64Equal("Direct ipv6 import updates received", 20, x.ImportUpdates.Received, t)
	assert.Int64Equal("Direct ipv6 import updates rejected", 21, x.ImportUpdates.Rejected, t)
	assert.Int64Equal("Direct ipv6 import updates filtered", 22, x.ImportUpdates.Filtered, t)
	assert.Int64Equal("Direct ipv6 import updates ignored", 23, x.ImportUpdates.Ignored, t)
	assert.Int64Equal("Direct ipv6 import updates accepted", 24, x.ImportUpdates.Accepted, t)
	assert.Int64Equal("Direct ipv6 import withdraws received", 25, x.ImportWithdraws.Received, t)
	assert.Int64Equal("Direct ipv6 import withdraws rejected", 26, x.ImportWithdraws.Rejected, t)
	assert.Int64Equal("Direct ipv6 import withdraws filtered", 27, x.ImportWithdraws.Filtered, t)
	assert.Int64Equal("Direct ipv6 import withdraws ignored", 28, x.ImportWithdraws.Ignored, t)
	assert.Int64Equal("Direct ipv6 import withdraws accepted", 29, x.ImportWithdraws.Accepted, t)
	assert.Int64Equal("Direct ipv6 export updates received", 30, x.ExportUpdates.Received, t)
	assert.Int64Equal("Direct ipv6 export updates rejected", 31, x.ExportUpdates.Rejected, t)
	assert.Int64Equal("Direct ipv6 export updates filtered", 32, x.ExportUpdates.Filtered, t)
	assert.Int64Equal("Direct ipv6 export updates ignored", 33, x.ExportUpdates.Ignored, t)
	assert.Int64Equal("Direct ipv6 export updates accepted", 34, x.ExportUpdates.Accepted, t)
	assert.Int64Equal("Direct ipv6 export withdraws received", 35, x.ExportWithdraws.Received, t)
	assert.Int64Equal("Direct ipv6 export withdraws rejected", 36, x.ExportWithdraws.Rejected, t)
	assert.Int64Equal("Direct ipv6 export withdraws filtered", 37, x.ExportWithdraws.Filtered, t)
	assert.Int64Equal("Direct ipv6 export withdraws ignored", 38, x.ExportWithdraws.Ignored, t)
	assert.Int64Equal("Direct ipv6 export withdraws accepted", 0, x.ExportWithdraws.Accepted, t)

	x = p[3]
	assert.StringEqual("OSPF ipv4 name", "ospf1", x.Name, t)
	assert.IntEqual("OSPF ipv4 proto", int(protocol.OSPF), int(x.Proto), t)
	assert.StringEqual("OSPF ipv4 ip version", "4", x.IPVersion, t)
	assert.Int64Equal("OSPF ipv4 imported", 4, x.Imported, t)
	assert.Int64Equal("OSPF ipv4 exported", 2, x.Exported, t)
	assert.Int64Equal("OSPF ipv4 filtered", 3, x.Filtered, t)
	assert.Int64Equal("OSPF ipv4 preferred", 1, x.Preferred, t)
}

func TestOSPFOldTimeFormat(t *testing.T) {
	data := "ospf1    OSPF      master   up     1481973060  Running\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "4")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "ospf1", x.Name, t)
	assert.IntEqual("proto", int(protocol.OSPF), int(x.Proto), t)
	assert.IntEqual("up", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("preferred", 100, x.Preferred, t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
}

func TestOSPFCurrentTimeFormat(t *testing.T) {
	data := "ospf1    OSPF      master   up     00:01:00  Running\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := ParseProtocols([]byte(data), "4")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "ospf1", x.Name, t)
	assert.IntEqual("proto", int(protocol.OSPF), int(x.Proto), t)
	assert.IntEqual("up", 1, x.Up, t)
	assert.Int64Equal("imported", 12, x.Imported, t)
	assert.Int64Equal("exported", 34, x.Exported, t)
	assert.Int64Equal("preferred", 100, x.Preferred, t)
	assert.StringEqual("ipVersion", "4", x.IPVersion, t)
	assert.IntEqual("uptime", 60, x.Uptime, t)
}

func TestRPKIUp(t *testing.T) {
	data := "rpki1      RPKI       ---        up     2021-12-31 13:04:29  Established"
	p := ParseProtocols([]byte(data), "4")
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "rpki1", x.Name, t)
	assert.IntEqual("proto", int(protocol.RPKI), int(x.Proto), t)
	assert.StringEqual("state", "Established", x.State, t)
	assert.IntEqual("up", 1, x.Up, t)
}
