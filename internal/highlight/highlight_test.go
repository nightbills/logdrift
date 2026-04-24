package highlight_test

import (
	"strings"
	"testing"

	"github.com/user/logdrift/internal/highlight"
)

func TestApply_NoKeywords(t *testing.T) {
	line := "hello world"
	got := highlight.Apply(line, highlight.Options{})
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_SingleKeyword(t *testing.T) {
	line := "error occurred in service"
	got := highlight.Apply(line, highlight.Options{Keywords: []string{"error"}})
	if !strings.Contains(got, "error") {
		t.Error("expected keyword to still appear in output")
	}
	if !strings.Contains(got, "\033[") {
		t.Error("expected ANSI escape codes in output")
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	line := "ERROR: something went wrong"
	got := highlight.Apply(line, highlight.Options{Keywords: []string{"error"}})
	if !strings.Contains(got, "\033[") {
		t.Error("expected ANSI codes for case-insensitive match")
	}
}

func TestApply_CaseSensitive_NoMatch(t *testing.T) {
	line := "ERROR: something went wrong"
	got := highlight.Apply(line, highlight.Options{
		Keywords:      []string{"error"},
		CaseSensitive: true,
	})
	if strings.Contains(got, "\033[") {
		t.Error("expected no ANSI codes for case-sensitive non-match")
	}
}

func TestApply_MultipleKeywords(t *testing.T) {
	line := "warn: disk usage high on node-1"
	got := highlight.Apply(line, highlight.Options{
		Keywords: []string{"warn", "disk"},
	})
	count := strings.Count(got, "\033[")
	// Each keyword match adds 2 escape sequences (open + reset)
	if count < 4 {
		t.Errorf("expected at least 4 ANSI codes, got %d", count)
	}
}

func TestApply_EmptyKeywordSkipped(t *testing.T) {
	line := "some log line"
	got := highlight.Apply(line, highlight.Options{Keywords: []string{"", "log"}})
	if !strings.Contains(got, "\033[") {
		t.Error("expected ANSI codes for non-empty keyword")
	}
}

func TestApply_KeywordNotPresent(t *testing.T) {
	line := "everything is fine"
	got := highlight.Apply(line, highlight.Options{Keywords: []string{"error"}})
	if strings.Contains(got, "\033[") {
		t.Error("expected no ANSI codes when keyword absent")
	}
}
