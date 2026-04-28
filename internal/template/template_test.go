package template

import (
	"testing"
)

func TestNew_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := New("{{.unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestNew_ValidTemplate_NoError(t *testing.T) {
	_, err := New("{{.level}} {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRender_NonJSON_ReturnedAsIs(t *testing.T) {
	r, _ := New("{{.level}} {{.msg}}")
	input := "plain text log line"
	if got := r.Render(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestRender_BasicFields(t *testing.T) {
	r, _ := New("[{{.level}}] {{.msg}}")
	line := `{"level":"info","msg":"hello world"}`
	want := "[info] hello world"
	if got := r.Render(line); got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestRender_MissingKey_EmptyString(t *testing.T) {
	r, _ := New("{{.level}} {{.missing}}")
	line := `{"level":"warn"}`
	want := "warn "
	if got := r.Render(line); got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestRender_NestedObject_JSONEncoded(t *testing.T) {
	r, _ := New("{{.meta}}")
	line := `{"meta":{"host":"srv1","pid":42}}`
	got := r.Render(line)
	if got == line {
		t.Fatal("expected nested object to be rendered, not raw line")
	}
	if len(got) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestRender_NumericValue_AsString(t *testing.T) {
	r, _ := New("code={{.code}}")
	line := `{"code":404}`
	want := "code=404"
	if got := r.Render(line); got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestRender_TrailingNewlineStripped(t *testing.T) {
	r, _ := New("{{.msg}}")
	line := "{\"msg\":\"hi\"}\n"
	want := "hi"
	if got := r.Render(line); got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
