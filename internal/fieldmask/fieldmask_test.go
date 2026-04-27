package fieldmask_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/fieldmask"
)

func jsonKeys(t *testing.T, s string) map[string]struct{} {
	t.Helper()
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	keys := make(map[string]struct{}, len(obj))
	for k := range obj {
		keys[k] = struct{}{}
	}
	return keys
}

const sampleLine = `{"level":"info","service":"api","msg":"started","ts":"2024-01-01T00:00:00Z"}`

func TestApply_NonJSON_ReturnedAsIs(t *testing.T) {
	m := fieldmask.New("", fieldmask.ModeExclude)
	got := m.Apply("not json at all")
	if got != "not json at all" {
		t.Errorf("expected raw line back, got %q", got)
	}
}

func TestApply_ExcludeEmpty_AllFieldsPresent(t *testing.T) {
	m := fieldmask.New("", fieldmask.ModeExclude)
	got := m.Apply(sampleLine)
	keys := jsonKeys(t, got)
	for _, k := range []string{"level", "service", "msg", "ts"} {
		if _, ok := keys[k]; !ok {
			t.Errorf("expected key %q to be present", k)
		}
	}
}

func TestApply_IncludeFields_OnlyListedPresent(t *testing.T) {
	m := fieldmask.New("level,msg", fieldmask.ModeInclude)
	got := m.Apply(sampleLine)
	keys := jsonKeys(t, got)
	if _, ok := keys["level"]; !ok {
		t.Error("expected 'level' to be present")
	}
	if _, ok := keys["msg"]; !ok {
		t.Error("expected 'msg' to be present")
	}
	if _, ok := keys["service"]; ok {
		t.Error("expected 'service' to be absent")
	}
	if _, ok := keys["ts"]; ok {
		t.Error("expected 'ts' to be absent")
	}
}

func TestApply_ExcludeFields_ListedAbsent(t *testing.T) {
	m := fieldmask.New("ts,service", fieldmask.ModeExclude)
	got := m.Apply(sampleLine)
	keys := jsonKeys(t, got)
	if _, ok := keys["ts"]; ok {
		t.Error("expected 'ts' to be excluded")
	}
	if _, ok := keys["service"]; ok {
		t.Error("expected 'service' to be excluded")
	}
	if _, ok := keys["level"]; !ok {
		t.Error("expected 'level' to remain")
	}
}

func TestApply_IncludeEmpty_NoFieldsReturned(t *testing.T) {
	m := fieldmask.New("", fieldmask.ModeInclude)
	got := m.Apply(sampleLine)
	keys := jsonKeys(t, got)
	if len(keys) != 0 {
		t.Errorf("expected empty object, got keys: %v", keys)
	}
}

func TestApply_WhitespaceInFieldList(t *testing.T) {
	m := fieldmask.New(" level , msg ", fieldmask.ModeInclude)
	got := m.Apply(sampleLine)
	keys := jsonKeys(t, got)
	if _, ok := keys["level"]; !ok {
		t.Error("expected 'level' after trimming whitespace")
	}
}
