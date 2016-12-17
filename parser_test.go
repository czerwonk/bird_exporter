package main

import "testing"

func TestEstablishedBirdOldFormat(t *testing.T) {
	data := "foo    BGP      master   up     1481973060  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	s := parseOutput([]byte(data), 4)
	assertInt("sessions", 1, len(s), t)

	x := s[0]
	assertString("name", "foo", x.name, t)
	assertInt("established", 1, x.established, t)
	assertInt("imported", 12, int(x.imported), t)
	assertInt("exported", 34, int(x.exported), t)
	assertInt("ipVersion", 4, x.ipVersion, t)
}

func TestEstablishedBirdCurrentFormat(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	s := parseOutput([]byte(data), 4)
	assertInt("sessions", 1, len(s), t)

	x := s[0]
	assertString("name", "foo", x.name, t)
	assertInt("established", 1, x.established, t)
	assertInt("imported", 12, int(x.imported), t)
	assertInt("exported", 34, int(x.exported), t)
	assertInt("ipVersion", 4, x.ipVersion, t)
	assertInt("uptime", 60, int(x.uptime), t)
}

func TestIpv6(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\nbar\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nxxx"
	s := parseOutput([]byte(data), 6)
	assertInt("sessions", 1, len(s), t)

	x := s[0]
	assertInt("ipVersion", 6, x.ipVersion, t)
}

func TestActive(t *testing.T) {
	data := "bar    BGP      master   start   2016-01-01    Active\ntest\nbar"
	s := parseOutput([]byte(data), 4)
	assertInt("sessions", 1, len(s), t)

	x := s[0]
	assertString("name", "bar", x.name, t)
	assertInt("established", 0, x.established, t)
	assertInt("imported", 0, int(x.imported), t)
	assertInt("exported", 0, int(x.exported), t)
	assertInt("ipVersion", 4, x.ipVersion, t)
	assertInt("uptime", 0, int(x.uptime), t)
}

func Test2Sessions(t *testing.T) {
	data := "foo    BGP      master   up     00:01:00  Established\ntest\n  Routes:         12 imported, 1 filtered, 34 exported, 100 preferred\nbar    BGP      master   start   2016-01-01    Active\nxxx"
	s := parseOutput([]byte(data), 4)
	assertInt("sessions", 2, len(s), t)
}

func assertString(name string, expected string, current string, t *testing.T) {
	if current != expected {
		t.Fatalf("%s: expected %s but got %s", name, expected, current)
	}
}

func assertInt(name string, expected int, current int, t *testing.T) {
	if current != expected {
		t.Fatalf("%s: expected %d but got %d", name, expected, current)
	}
}
