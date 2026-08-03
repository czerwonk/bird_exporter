package parser

import (
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	sample := `
BIRD 2.17.1 ready.
BIRD 2.17.1
Router ID is 10.72.0.164
Hostname is store-blob-pp014
Current server time is 2025-08-21 10:27:57.219
Last reboot on 2025-07-31 14:23:03.232
Last reconfiguration on 2025-08-07 09:53:19.531
Daemon is up and running
`

	s := ParseStatus([]byte(sample))

	if s.Version != "2.17.1" {
		t.Errorf("expected Version=2.17.1, got %q", s.Version)
	}
	if s.RouterID != "10.72.0.164" {
		t.Errorf("expected RouterID=10.72.0.164, got %q", s.RouterID)
	}

	expectedServerTime := time.Date(2025, 8, 21, 10, 27, 57, 219000000, time.UTC)
	if !s.ServerTime.Equal(expectedServerTime) {
		t.Errorf("expected ServerTime=%v, got %v", expectedServerTime, s.ServerTime)
	}

	expectedLastReboot := time.Date(2025, 7, 31, 14, 23, 3, 232000000, time.UTC)
	if !s.LastReboot.Equal(expectedLastReboot) {
		t.Errorf("expected LastReboot=%v, got %v", expectedLastReboot, s.LastReboot)
	}

	expectedLastReconfig := time.Date(2025, 8, 7, 9, 53, 19, 531000000, time.UTC)
	if !s.LastReconfig.Equal(expectedLastReconfig) {
		t.Errorf("expected LastReconfig=%v, got %v", expectedLastReconfig, s.LastReconfig)
	}

	if !s.DaemonUp {
		t.Errorf("expected DaemonUp=true, got false")
	}
}

func TestParseStatusMissingFields(t *testing.T) {
	sample := `
BIRD 2.17.1
Daemon is down
`

	s := ParseStatus([]byte(sample))

	if s.Version != "2.17.1" {
		t.Errorf("expected Version=2.17.1, got %q", s.Version)
	}
	if s.DaemonUp {
		t.Errorf("expected DaemonUp=false, got true")
	}

	// Other fields should be zero-values
	if !s.ServerTime.IsZero() {
		t.Errorf("expected ServerTime zero, got %v", s.ServerTime)
	}
	if !s.LastReboot.IsZero() {
		t.Errorf("expected LastReboot zero, got %v", s.LastReboot)
	}
	if !s.LastReconfig.IsZero() {
		t.Errorf("expected LastReconfig zero, got %v", s.LastReconfig)
	}
	if s.RouterID != "" {
		t.Errorf("expected RouterID empty, got %q", s.RouterID)
	}
}

func TestParseStatusMalformedTime(t *testing.T) {
	sample := `
BIRD 2.17.1
Current server time is not-a-time
Daemon is up
`

	s := ParseStatus([]byte(sample))

	if !s.DaemonUp {
		t.Errorf("expected DaemonUp=true, got false")
	}
	if !s.ServerTime.IsZero() {
		t.Errorf("expected ServerTime zero due to parse error, got %v", s.ServerTime)
	}
}

func TestParseStatusWithErrorRejectsMalformedFields(t *testing.T) {
	tests := []struct {
		name   string
		sample string
	}{
		{name: "missing daemon state", sample: "BIRD 2.17.1\n"},
		{name: "unknown daemon state", sample: "BIRD 2.17.1\nDaemon is sideways\n"},
		{name: "malformed server time", sample: "BIRD 2.17.1\nCurrent server time is not-a-time\nDaemon is up\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseStatusWithError([]byte(tt.sample)); err == nil {
				t.Fatal("expected malformed BIRD status to return an error")
			}
		})
	}
}

func TestParseStatusWithErrorAcceptsDaemonDown(t *testing.T) {
	s, err := ParseStatusWithError([]byte("BIRD 2.17.1\nDaemon is down\n"))
	if err != nil {
		t.Fatalf("expected daemon-down status to parse: %v", err)
	}
	if s.DaemonUp {
		t.Fatal("expected daemon-down status to remain down")
	}
}
