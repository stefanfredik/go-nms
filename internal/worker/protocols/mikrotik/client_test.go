package mikrotik_test

import (
	"testing"
	"time"

	mikrotik "github.com/yourorg/nms-go/internal/worker/protocols/mikrotik"
)

func TestParseRouterOSUptime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		// Legacy CLI format (XwXdXhXmXs)
		{"0s", 0},
		{"30s", 30 * time.Second},
		{"5m30s", 5*time.Minute + 30*time.Second},
		{"3h25m10s", 3*time.Hour + 25*time.Minute + 10*time.Second},
		{"2d3h", 2*24*time.Hour + 3*time.Hour},
		{"1w2d3h4m5s", 7*24*time.Hour + 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second},
		{"14h55m", 14*time.Hour + 55*time.Minute},

		// RouterOS API format (DDdHH:MM:SS)
		{"20:47:30", 20*time.Hour + 47*time.Minute + 30*time.Second},
		{"3d14:25:10", 3*24*time.Hour + 14*time.Hour + 25*time.Minute + 10*time.Second},
		{"125d20:47:30", 125*24*time.Hour + 20*time.Hour + 47*time.Minute + 30*time.Second},
		{"1w3d08:00:00", 7*24*time.Hour + 3*24*time.Hour + 8*time.Hour},

		// Edge cases
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mikrotik.ParseRouterOSUptime(tt.input)
			if got != tt.expected {
				t.Errorf("ParseRouterOSUptime(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
