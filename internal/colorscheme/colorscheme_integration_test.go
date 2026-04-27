package colorscheme_test

import (
	"strings"
	"testing"

	"github.com/user/logdrift/internal/colorscheme"
)

// TestAllSchemes_LevelColorsApply verifies that every built-in scheme can
// wrap a sample message without panicking and that plain produces no escapes.
func TestAllSchemes_LevelColorsApply(t *testing.T) {
	levels := []struct {
		name  string
		color func(colorscheme.Scheme) string
	}{
		{"error", func(s colorscheme.Scheme) string { return s.Error }},
		{"warn", func(s colorscheme.Scheme) string { return s.Warn }},
		{"info", func(s colorscheme.Scheme) string { return s.Info }},
		{"debug", func(s colorscheme.Scheme) string { return s.Debug }},
	}

	for _, schemeName := range colorscheme.Names() {
		s := colorscheme.Get(schemeName)
		for _, lv := range levels {
			out := s.Apply(lv.color(s), "msg")
			if !strings.Contains(out, "msg") {
				t.Errorf("scheme=%s level=%s: output missing original text", schemeName, lv.name)
			}
			if schemeName == "plain" && strings.Contains(out, "\033[") {
				t.Errorf("scheme=plain level=%s: unexpected escape in output", lv.name)
			}
		}
	}
}

// TestScheme_ServiceAndKeyColors ensures structural fields are colored.
func TestScheme_ServiceAndKeyColors(t *testing.T) {
	s := colorscheme.Get("default")
	service := s.Apply(s.Service, "auth-svc")
	key := s.Apply(s.Key, "request_id")

	if !strings.HasPrefix(service, s.Service) {
		t.Error("service color not applied")
	}
	if !strings.HasPrefix(key, s.Key) {
		t.Error("key color not applied")
	}
}
