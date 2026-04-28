package retag_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/retag"
)

func jsonKeys(t *testing.T, line string) map[string]json.RawMessage {
	t.Helper()
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return obj
}

func TestApply_NonJSON_ReturnedAsIs(t *testing.T) {
	r := retag.New(map[string]string{"old": "new"})
	input := "not json at all"
	if got := r.Apply(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestApply_EmptyMapping_NoChange(t *testing.T) {
	r := retag.New(map[string]string{})
	input := `{"level":"info","msg":"hello"}`
	if got := r.Apply(input); got != input {
		t.Errorf("expected input unchanged, got %q", got)
	}
}

func TestApply_RenamesField(t *testing.T) {
	r := retag.New(map[string]string{"message": "msg"})
	input := `{"level":"info","message":"hello world"}`
	got := r.Apply(input)
	keys := jsonKeys(t, got)
	if _, ok := keys["message"]; ok {
		t.Error("old key 'message' should not be present")
	}
	if _, ok := keys["msg"]; !ok {
		t.Error("new key 'msg' should be present")
	}
}

func TestApply_MissingSourceField_NoError(t *testing.T) {
	r := retag.New(map[string]string{"nonexistent": "renamed"})
	input := `{"level":"info","msg":"hello"}`
	got := r.Apply(input)
	keys := jsonKeys(t, got)
	if _, ok := keys["renamed"]; ok {
		t.Error("renamed key should not appear when source is absent")
	}
	if _, ok := keys["level"]; !ok {
		t.Error("unrelated key 'level' should still be present")
	}
}

func TestApply_MultipleRenames(t *testing.T) {
	r := retag.New(map[string]string{
		"svc": "service",
		"ts":  "time",
	})
	input := `{"svc":"auth","ts":"2024-01-01T00:00:00Z","msg":"ok"}`
	got := r.Apply(input)
	keys := jsonKeys(t, got)
	for _, old := range []string{"svc", "ts"} {
		if _, ok := keys[old]; ok {
			t.Errorf("old key %q should not be present after rename", old)
		}
	}
	for _, newKey := range []string{"service", "time"} {
		if _, ok := keys[newKey]; !ok {
			t.Errorf("new key %q should be present after rename", newKey)
		}
	}
}

func TestApply_EmptyKeyInMapping_Ignored(t *testing.T) {
	r := retag.New(map[string]string{"": "newkey", "level": ""})
	input := `{"level":"info","msg":"hello"}`
	got := r.Apply(input)
	// level should still be present since the mapping entry has empty newKey
	keys := jsonKeys(t, got)
	if _, ok := keys["level"]; !ok {
		t.Error("key 'level' should remain unchanged when mapping target is empty")
	}
}
