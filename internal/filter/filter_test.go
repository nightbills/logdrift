package filter

import (
	"testing"
)

func TestParse_ValidJSON(t *testing.T) {
	line := `{"level":"info","msg":"started","service":"api"}`
	entry := Parse(line)
	if entry == nil {
		t.Fatal("expected non-nil entry for valid JSON")
	}
	if entry["level"] != "info" {
		t.Errorf("expected level=info, got %v", entry["level"])
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	if entry := Parse("not json at all"); entry != nil {
		t.Error("expected nil for non-JSON line")
	}
	if entry := Parse(""); entry != nil {
		t.Error("expected nil for empty line")
	}
}

func TestMatch_NoFilters(t *testing.T) {
	raw := `{"level":"debug","msg":"ping"}`
	entry := Parse(raw)
	if !Match(entry, raw, Options{}) {
		t.Error("expected match with empty options")
	}
}

func TestMatch_LevelFilter(t *testing.T) {
	raw := `{"level":"error","msg":"boom","service":"worker"}`
	entry := Parse(raw)

	if !Match(entry, raw, Options{Level: "error"}) {
		t.Error("expected match on level=error")
	}
	if Match(entry, raw, Options{Level: "info"}) {
		t.Error("expected no match on level=info")
	}
}

func TestMatch_ServiceFilter(t *testing.T) {
	raw := `{"level":"info","service":"gateway","msg":"request"}`
	entry := Parse(raw)

	if !Match(entry, raw, Options{Service: "gateway"}) {
		t.Error("expected match on service=gateway")
	}
	if Match(entry, raw, Options{Service: "worker"}) {
		t.Error("expected no match on service=worker")
	}
}

func TestMatch_ContainsFilter(t *testing.T) {
	raw := `{"level":"warn","msg":"disk space low"}`
	entry := Parse(raw)

	if !Match(entry, raw, Options{Contains: "disk space"}) {
		t.Error("expected match when substring present")
	}
	if Match(entry, raw, Options{Contains: "cpu"}) {
		t.Error("expected no match when substring absent")
	}
}

func TestMatch_CombinedFilters(t *testing.T) {
	raw := `{"level":"error","service":"api","msg":"timeout"}`
	entry := Parse(raw)
	opts := Options{Level: "error", Service: "api", Contains: "timeout"}

	if !Match(entry, raw, opts) {
		t.Error("expected match with all filters satisfied")
	}

	opts.Service = "db"
	if Match(entry, raw, opts) {
		t.Error("expected no match when service differs")
	}
}
