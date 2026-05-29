package toptags

import (
	"testing"
	"time"
)

func TestResolveCutoff(t *testing.T) {
	today := time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		window  string
		want    time.Time
		wantErr bool
	}{
		{name: "year", window: "1y", want: time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC)},
		{name: "month", window: "6m", want: time.Date(2025, 11, 29, 0, 0, 0, 0, time.UTC)},
		{name: "day", window: "90d", want: time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)},
		{name: "composite", window: "1y2m10d", want: time.Date(2025, 3, 19, 0, 0, 0, 0, time.UTC)},
		{name: "invalid unit", window: "1w", wantErr: true},
		{name: "missing unit", window: "10", wantErr: true},
		{name: "empty", window: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveCutoff(today, tt.window)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("resolveCutoff(%q) error = nil, want error", tt.window)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveCutoff(%q) error = %v", tt.window, err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("resolveCutoff(%q) = %s, want %s", tt.window, got.Format("2006-01-02"), tt.want.Format("2006-01-02"))
			}
		})
	}
}
