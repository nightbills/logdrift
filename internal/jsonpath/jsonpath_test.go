package jsonpath_test

import (
	"testing"

	"github.com/user/logdrift/internal/jsonpath"
)

func TestExtract_TopLevelString(t *testing.T) {
	val, err := jsonpath.Extract(`{"level":"info"}`, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "info" {
		t.Errorf("expected \"info\", got %q", val)
	}
}

func TestExtract_NestedKey(t *testing.T) {
	val, err := jsonpath.Extract(`{"meta":{"env":"prod"}}`, "meta.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "prod" {
		t.Errorf("expected \"prod\", got %q", val)
	}
}

func TestExtract_DeeplyNested(t *testing.T) {
	line := `{"a":{"b":{"c":"deep"}}}`
	val, err := jsonpath.Extract(line, "a.b.c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "deep" {
		t.Errorf("expected \"deep\", got %q", val)
	}
}

func TestExtract_NumericValue(t *testing.T) {
	val, err := jsonpath.Extract(`{"latency":42}`, "latency")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "42" {
		t.Errorf("expected \"42\", got %q", val)
	}
}

func TestExtract_BoolValue(t *testing.T) {
	val, err := jsonpath.Extract(`{"ok":true}`, "ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "true" {
		t.Errorf("expected \"true\", got %q", val)
	}
}

func TestExtract_MissingKey(t *testing.T) {
	_, err := jsonpath.Extract(`{"level":"info"}`, "missing")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestExtract_IntermediateNotObject(t *testing.T) {
	_, err := jsonpath.Extract(`{"meta":"string"}`, "meta.env")
	if err == nil {
		t.Fatal("expected error when intermediate key is not an object")
	}
}

func TestExtract_InvalidJSON(t *testing.T) {
	_, err := jsonpath.Extract(`not json`, "level")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestExtract_NullValue(t *testing.T) {
	val, err := jsonpath.Extract(`{"field":null}`, "field")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for null, got %q", val)
	}
}
