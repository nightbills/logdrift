package fieldmask_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/fieldmask"
)

// TestRoundTrip ensures that values are preserved faithfully after masking.
func TestRoundTrip_ValuesPreserved(t *testing.T) {
	line := `{"level":"error","code":42,"ok":false,"meta":{"host":"box1"}}`

	m := fieldmask.New("level,code,ok,meta", fieldmask.ModeInclude)
	out := m.Apply(line)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	cases := map[string]string{
		"level": `"error"`,
		"code":  `42`,
		"ok":    `false`,
		"meta":  `{"host":"box1"}`,
	}
	for key, want := range cases {
		got, ok := obj[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if string(got) != want {
			t.Errorf("key %q: got %s, want %s", key, got, want)
		}
	}
}

// TestChainedMasks simulates applying an exclude mask followed by an include mask.
func TestChainedMasks(t *testing.T) {
	line := `{"level":"info","ts":"2024-01-01","service":"svc","msg":"hello"}`

	// First strip ts
	exclude := fieldmask.New("ts", fieldmask.ModeExclude)
	intermediate := exclude.Apply(line)

	// Then keep only level and msg
	include := fieldmask.New("level,msg", fieldmask.ModeInclude)
	out := include.Apply(intermediate)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("chained output is not valid JSON: %v", err)
	}
	if len(obj) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(obj), obj)
	}
	if _, ok := obj["ts"]; ok {
		t.Error("ts should have been excluded")
	}
	if _, ok := obj["service"]; ok {
		t.Error("service should have been excluded by include mask")
	}
}
