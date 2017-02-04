package main

import (
	"testing"

	"github.com/czerwonk/testutils/assert"
)

func TestEstablishedBirdOldFormat(t *testing.T) {
	data := "foo    BGP      master   up     1481973060  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.name, t)
	assert.IntEqual("established", 1, x.established, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
}

func TestEstablishedBirdCurrentFormat(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "foo", x.name, t)
	assert.IntEqual("established", 1, x.established, t)
	assert.Int64Equal("imported", 12, x.imported, t)
	assert.Int64Equal("exported", 34, x.exported, t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
	assert.IntEqual("uptime", 60, x.uptime, t)
}

func TestIpv6(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	p := parseOutput([]byte(data), 6)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.IntEqual("ipVersion", 6, x.ipVersion, t)
}

func TestActive(t *testing.T) {
	data := "bar    BGP      master   start   2016-01-01    Active\ntest\nbar"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 1, len(p), t)

	x := p[0]
	assert.StringEqual("name", "bar", x.name, t)
	assert.IntEqual("established", 0, x.established, t)
	assert.IntEqual("imported", 0, int(x.imported), t)
	assert.IntEqual("exported", 0, int(x.exported), t)
	assert.IntEqual("ipVersion", 4, x.ipVersion, t)
	assert.IntEqual("uptime", 0, int(x.uptime), t)
}

func Test2Sessions(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nbar    BGP      master   start   2016-01-01    Active\nxxx"
	p := parseOutput([]byte(data), 4)
	assert.IntEqual("protocols", 2, len(p), t)
}
