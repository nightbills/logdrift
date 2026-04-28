package retag_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/retag"
)

// TestChainedRetaggers verifies that applying two Retaggers in sequence
// produces the expected cumulative rename result.
func TestChainedRetaggers(t *testing.T) {
	first := retag.New(map[string]string{"message": "msg"})
	second := retag.New(map[string]string{"svc": "service"})

	input := `{"svc":"payments","message":"charge accepted","level":"info"}`
	intermediate := first.Apply(input)
	final := second.Apply(intermediate)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(final), &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	expected := map[string]bool{
		"service": true,
		"msg":     true,
		"level":   true,
	}
	for key := range expected {
		if _, ok := obj[key]; !ok {
			t.Errorf("expected key %q to be present in final output", key)
		}
	}
	absent := []string{"message", "svc"}
	for _, key := range absent {
		if _, ok := obj[key]; ok {
			t.Errorf("key %q should have been renamed away", key)
		}
	}
}

// TestValuePreservedAfterRename ensures field values survive the rename.
func TestValuePreservedAfterRename(t *testing.T) {
	r := retag.New(map[string]string{"message": "msg"})
	input := `{"message":"hello world","level":"debug"}`
	out := r.Apply(input)

	var obj map[string]string
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["msg"] != "hello world" {
		t.Errorf("expected msg=\"hello world\", got %q", obj["msg"])
	}
}
