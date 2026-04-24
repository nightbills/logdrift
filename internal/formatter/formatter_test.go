package formatter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logdrift/internal/formatter"
)

var defaultOpts = formatter.Options{NoColor: true}

func TestFormat_NonJSON(t *testing.T) {
	raw := "not json at all"
	out := formatter.Format(raw, defaultOpts)
	if out != raw {
		t.Errorf("expected raw line passthrough, got %q", out)
	}
}

func TestFormat_BasicFields(t *testing.T) {
	line := `{"level":"info","msg":"server started","service":"api"}`
	out := formatter.Format(line, defaultOpts)

	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected level in output, got: %s", out)
	}
	if !strings.Contains(out, "server started") {
		t.Errorf("expected msg in output, got: %s", out)
	}
	if !strings.Contains(out, "(api)") {
		t.Errorf("expected service in output, got: %s", out)
	}
}

func TestFormat_ErrorLevel(t *testing.T) {
	line := `{"level":"error","msg":"connection refused"}`
	out := formatter.Format(line, defaultOpts)
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", out)
	}
}

func TestFormat_WarnLevel(t *testing.T) {
	line := `{"level":"warn","msg":"retrying"}`
	out := formatter.Format(line, defaultOpts)
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
}

func TestFormat_TimeParsed(t *testing.T) {
	line := `{"level":"info","time":"2024-03-15T14:05:00Z","msg":"ok"}`
	out := formatter.Format(line, formatter.Options{NoColor: true, TimeFormat: "15:04:05"})
	if !strings.Contains(out, "14:05:00") {
		t.Errorf("expected formatted time in output, got: %s", out)
	}
}

func TestFormat_AlternativeMessageKey(t *testing.T) {
	line := `{"level":"info","message":"hello world"}`
	out := formatter.Format(line, defaultOpts)
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestFormat_MissingLevel(t *testing.T) {
	line := `{"msg":"no level here"}`
	out := formatter.Format(line, defaultOpts)
	if !strings.Contains(out, "[???]") {
		t.Errorf("expected fallback level [???], got: %s", out)
	}
}

func TestFormat_ColorEnabled(t *testing.T) {
	line := `{"level":"error","msg":"boom"}`
	out := formatter.Format(line, formatter.Options{NoColor: false})
	// ANSI escape sequences should be present
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes in colored output, got: %s", out)
	}
}
