package parser

import (
	"fmt"
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
)

func ParseStatus(b []byte) *Status {
	s, _ := ParseStatusWithError(b)
	return s
}

// ParseStatusWithError parses BIRD status output and reports malformed fields
// that would otherwise make daemon health look successful.
func ParseStatusWithError(b []byte) (*Status, error) {
	s := &Status{}
	out := string(b)
	var parseErr error

	if m := reVersion.FindStringSubmatch(out); len(m) > 1 {
		s.Version = m[1]
	}
	if m := reRouterID.FindStringSubmatch(out); len(m) > 1 {
		s.RouterID = m[1]
	}
	if m := reServerTime.FindStringSubmatch(out); len(m) > 1 {
		t, err := time.ParseInLocation(birdTimestampLayout, strings.TrimSpace(m[1]), time.Local)
		if err != nil {
			parseErr = fmt.Errorf("parse BIRD server time: %w", err)
		} else {
			s.ServerTime = t
		}
	}
	if m := reReboot.FindStringSubmatch(out); len(m) > 1 {
		t, err := time.ParseInLocation(birdTimestampLayout, strings.TrimSpace(m[1]), time.Local)
		if err != nil {
			if parseErr == nil {
				parseErr = fmt.Errorf("parse BIRD last reboot time: %w", err)
			}
		} else {
			s.LastReboot = t
		}
	}
	if m := reReconfig.FindStringSubmatch(out); len(m) > 1 {
		t, err := time.ParseInLocation(birdTimestampLayout, strings.TrimSpace(m[1]), time.Local)
		if err != nil {
			if parseErr == nil {
				parseErr = fmt.Errorf("parse BIRD last reconfiguration time: %w", err)
			}
		} else {
			s.LastReconfig = t
		}
	}

	m := reDaemon.FindStringSubmatch(out)
	if len(m) <= 1 {
		if parseErr == nil {
			parseErr = fmt.Errorf("BIRD status reply is missing daemon state")
		}
	} else {
		daemonState := strings.TrimSpace(m[1])
		switch {
		case daemonState == "up" || strings.HasPrefix(daemonState, "up "):
			s.DaemonUp = true
		case daemonState == "down" || strings.HasPrefix(daemonState, "down "):
			s.DaemonUp = false
		default:
			if parseErr == nil {
				parseErr = fmt.Errorf("unrecognized BIRD daemon state %q", daemonState)
			}
		}
	}

	return s, parseErr
}
