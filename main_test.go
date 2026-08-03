package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatusSocketPath(t *testing.T) {
	tests := []struct {
		name        string
		birdV2      bool
		birdEnabled bool
		want        string
	}{
		{name: "BIRD v2", birdV2: true, want: "/run/bird/bird.ctl"},
		{name: "BIRD v1 dual stack", birdEnabled: true, want: "/run/bird/bird.ctl"},
		{name: "BIRD v1 IPv6 only", want: "/run/bird/bird6.ctl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusSocketPath(tt.birdV2, tt.birdEnabled, "/run/bird/bird.ctl", "/run/bird/bird6.ctl")
			require.Equal(t, tt.want, got)
		})
	}
}
