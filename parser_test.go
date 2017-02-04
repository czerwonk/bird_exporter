package main

import (
	"testing"

	"github.com/czerwonk/testutils/assert"
)

func TestEstablishedBgpOldTimeFormat(t *testing.T) {
	data := "foo    BGP      master   up     1481973060  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.name, t)
	assert.IntEqual("proto", BGP, x.proto, t)
	assert.IntEqual("established", 1, x.up, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
}

func TestEstablishedBgpCurrentTimeFormat(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.name, t)
	assert.IntEqual("proto", BGP, x.proto, t)
	assert.IntEqual("established", 1, x.up, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
	assert.IntEqual("uptime", 60, x.uptime, t)
}

func TestIpv6Bgp(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 6)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.IntEqual("ipVersion", 6, x.ipVersion, t)
}

func TestActiveBgp(t *testing.T) {
	data := "bar    BGP      master   start   2016-01-01    Active\ntest\nbar"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "bar", x.name, t)
	assert.IntEqual("proto", BGP, x.proto, t)
	assert.IntEqual("established", 0, x.up, t)
	assert.IntEqual("imported", 0, int(x.imported), t)
	assert.IntEqual("exported", 0, int(x.exported), t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
	assert.IntEqual("uptime", 0, int(x.uptime), t)
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
	assert.StringEqual("name", "ospf1", x.name, t)
	assert.IntEqual("proto", OSPF, x.proto, t)
	assert.IntEqual("established", 1, x.up, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
}

func TestOspfCurrentTimeFormat(t *testing.T) {
	data := "ospf1    OSPF      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "ospf1", x.name, t)
	assert.IntEqual("proto", OSPF, x.proto, t)
	assert.IntEqual("established", 1, x.up, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
	assert.IntEqual("uptime", 60, x.uptime, t)
}
