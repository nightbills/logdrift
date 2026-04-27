package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/levelfilter"
)

func TestParse_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  levelfilter.Level
	}{
		{"debug", levelfilter.LevelDebug},
		{"INFO", levelfilter.LevelInfo},
		{"Warn", levelfilter.LevelWarn},
		{"ERROR", levelfilter.LevelError},
		{"fatal", levelfilter.LevelFatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := levelfilter.Parse(tc.input)
			if got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParse_UnknownLevel(t *testing.T) {
	if got := levelfilter.Parse("trace"); got != levelfilter.LevelUnknown {
		t.Errorf("expected LevelUnknown, got %v", got)
	}
}

func TestFilter_Allow_MinInfo(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelInfo)

	if f.Allow("debug") {
		t.Error("debug should be blocked when min=info")
	}
	if !f.Allow("info") {
		t.Error("info should pass when min=info")
	}
	if !f.Allow("warn") {
		t.Error("warn should pass when min=info")
	}
	if !f.Allow("error") {
		t.Error("error should pass when min=info")
	}
	if !f.Allow("fatal") {
		t.Error("fatal should pass when min=info")
	}
}

func TestFilter_Allow_UnknownMin(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelUnknown)
	for _, lvl := range []string{"debug", "info", "warn", "error", "fatal", "trace"} {
		if !f.Allow(lvl) {
			t.Errorf("expected %q to pass when min=unknown", lvl)
		}
	}
}

func TestFilter_Allow_UnparsedLevelPassesThrough(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError)
	// "trace" is not a known level — should pass through rather than be dropped.
	if !f.Allow("trace") {
		t.Error("unrecognised level should pass through regardless of min threshold")
	}
}

func TestLevel_String(t *testing.T) {
	if s := levelfilter.LevelWarn.String(); s != "warn" {
		t.Errorf("LevelWarn.String() = %q, want \"warn\"", s)
	}
	if s := levelfilter.LevelUnknown.String(); s != "unknown" {
		t.Errorf("LevelUnknown.String() = %q, want \"unknown\"", s)
	}
}
