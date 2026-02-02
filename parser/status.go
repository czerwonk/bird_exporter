package parser

import (
	"regexp"
	"strings"
	"time"
)

type Status struct {
	Version      string
	RouterID     string
	ServerTime   time.Time
	LastReboot   time.Time
	LastReconfig time.Time
	DaemonUp     bool
}

var (
	reVersion    = regexp.MustCompile(`BIRD (\S+)`)
	reRouterID   = regexp.MustCompile(`Router ID is (\S+)`)
	reServerTime = regexp.MustCompile(`Current server time is (.+)`)
	reReboot     = regexp.MustCompile(`Last reboot on (.+)`)
	reReconfig   = regexp.MustCompile(`Last reconfiguration on (.+)`)
	reDaemon     = regexp.MustCompile(`Daemon is (.+)`)
	timeLayout   = "2006-01-02 15:04:05.000" // ref: https://pkg.go.dev/time#Parse
)

func ParseStatus(b []byte) *Status {
	s := &Status{}
	out := string(b)

	if m := reVersion.FindStringSubmatch(out); len(m) > 1 {
		s.Version = m[1]
	}
	if m := reRouterID.FindStringSubmatch(out); len(m) > 1 {
		s.RouterID = m[1]
	}
	if m := reServerTime.FindStringSubmatch(out); len(m) > 1 {
		if t, err := time.Parse(timeLayout, strings.TrimSpace(m[1])); err == nil {
			s.ServerTime = t
		}
	}
	if m := reReboot.FindStringSubmatch(out); len(m) > 1 {
		if t, err := time.Parse(timeLayout, strings.TrimSpace(m[1])); err == nil {
			s.LastReboot = t
		}
	}
	if m := reReconfig.FindStringSubmatch(out); len(m) > 1 {
		if t, err := time.Parse(timeLayout, strings.TrimSpace(m[1])); err == nil {
			s.LastReconfig = t
		}
	}
	if m := reDaemon.FindStringSubmatch(out); len(m) > 1 {
		s.DaemonUp = strings.Contains(m[1], "up")
	}

	return s
}
