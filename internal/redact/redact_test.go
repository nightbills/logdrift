package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/redact"
)

func TestApply_NonJSON_ReturnedAsIs(t *testing.T) {
	r := redact.New([]string{"password"}, "")
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_NoFields_ReturnedAsIs(t *testing.T) {
	r := redact.New(nil, "")
	line := `{"password":"secret","user":"alice"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_RedactsNamedField(t *testing.T) {
	r := redact.New([]string{"password"}, "")
	line := `{"user":"alice","password":"s3cr3t"}`
	got := r.Apply(line)

	var obj map[string]string
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Errorf("user field should be unchanged, got %q", obj["user"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	r := redact.New([]string{"token"}, "***")
	line := `{"token":"abc123"}`
	got := r.Apply(line)

	var obj map[string]string
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["token"] != "***" {
		t.Errorf("expected ***, got %q", obj["token"])
	}
}

func TestApply_CaseInsensitiveFieldMatch(t *testing.T) {
	r := redact.New([]string{"Authorization"}, "")
	line := `{"authorization":"Bearer xyz"}`
	got := r.Apply(line)

	var obj map[string]string
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["authorization"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", obj["authorization"])
	}
}

func TestApply_MultipleFields(t *testing.T) {
	r := redact.New([]string{"password", "ssn", "token"}, "")
	line := `{"user":"bob","password":"pass","ssn":"123-45-6789","token":"tok"}`
	got := r.Apply(line)

	var obj map[string]string
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	for _, field := range []string{"password", "ssn", "token"} {
		if obj[field] != "[REDACTED]" {
			t.Errorf("field %q: expected [REDACTED], got %q", field, obj[field])
		}
	}
	if obj["user"] != "bob" {
		t.Errorf("user should be unchanged, got %q", obj["user"])
	}
}
