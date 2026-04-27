package truncate

import (
	"strings"
	"testing"
)

func TestNew_ZeroWidth_NoTruncation(t *testing.T) {
	tr := New(0)
	long := strings.Repeat("a", 500)
	if got := tr.Line(long); got != long {
		t.Fatalf("expected unchanged line, got length %d", len(got))
	}
}

func TestLine_ShortLine_Unchanged(t *testing.T) {
	tr := New(80)
	input := "hello world"
	if got := tr.Line(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestLine_ExactWidth_Unchanged(t *testing.T) {
	tr := New(10)
	input := strings.Repeat("x", 10)
	if got := tr.Line(input); got != input {
		t.Fatalf("expected unchanged at exact width, got %q", got)
	}
}

func TestLine_ExceedsWidth_Truncated(t *testing.T) {
	tr := New(10)
	input := strings.Repeat("a", 20)
	got := tr.Line(input)
	if !strings.HasSuffix(got, ellipsis) {
		t.Fatalf("expected ellipsis suffix, got %q", got)
	}
	// visible chars before ellipsis should equal maxWidth
	visible := strings.TrimSuffix(got, ellipsis)
	if len(visible) != 10 {
		t.Fatalf("expected 10 visible chars before ellipsis, got %d", len(visible))
	}
}

func TestLine_ANSIEscapes_NotCounted(t *testing.T) {
	tr := New(5)
	// red colour code + 5 chars + reset — should NOT be truncated
	input := "\x1b[31m" + "hello" + "\x1b[0m"
	got := tr.Line(input)
	if strings.Contains(got, ellipsis) {
		t.Fatalf("ANSI sequences should not count toward width, got %q", got)
	}
}

func TestLine_ANSIEscapes_TruncatedAfterLimit(t *testing.T) {
	tr := New(3)
	input := "\x1b[32m" + "hello" + "\x1b[0m"
	got := tr.Line(input)
	if !strings.Contains(got, ellipsis) {
		t.Fatalf("expected truncation for 5 visible chars with limit 3, got %q", got)
	}
}

func TestLine_Unicode_CountedCorrectly(t *testing.T) {
	tr := New(4)
	// Each rune is one visible character regardless of byte width
	input := "日本語テスト"
	got := tr.Line(input)
	if !strings.HasSuffix(got, ellipsis) {
		t.Fatalf("expected truncation of unicode string, got %q", got)
	}
}
