package colorscheme_test

import (
	"strings"
	"testing"

	"github.com/user/logdrift/internal/colorscheme"
)

func TestGet_KnownScheme_ReturnsScheme(t *testing.T) {
	s := colorscheme.Get("dark")
	if s.Name != "dark" {
		t.Fatalf("expected dark, got %q", s.Name)
	}
}

func TestGet_UnknownScheme_ReturnsDefault(t *testing.T) {
	s := colorscheme.Get("nonexistent")
	if s.Name != "default" {
		t.Fatalf("expected default fallback, got %q", s.Name)
	}
}

func TestGet_PlainScheme_NoColors(t *testing.T) {
	s := colorscheme.Get("plain")
	if s.Error != "" || s.Info != "" {
		t.Fatal("plain scheme should have empty color codes")
	}
}

func TestApply_WrapsWithColorAndReset(t *testing.T) {
	s := colorscheme.Get("default")
	out := s.Apply(s.Error, "oops")
	if !strings.Contains(out, "oops") {
		t.Fatal("output should contain original text")
	}
	if !strings.Contains(out, colorscheme.Reset) {
		t.Fatal("output should contain reset sequence")
	}
	if !strings.HasPrefix(out, s.Error) {
		t.Fatal("output should start with color code")
	}
}

func TestApply_EmptyColor_ReturnsTextUnchanged(t *testing.T) {
	s := colorscheme.Get("plain")
	out := s.Apply(s.Info, "hello")
	if out != "hello" {
		t.Fatalf("expected unchanged text, got %q", out)
	}
}

func TestNames_ContainsBuiltins(t *testing.T) {
	names := colorscheme.Names()
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	for _, want := range []string{"default", "dark", "plain"} {
		if !set[want] {
			t.Errorf("expected scheme %q in Names()", want)
		}
	}
}
