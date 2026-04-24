package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logdrift/internal/output"
)

func TestNew_StdoutOnly(t *testing.T) {
	w, err := output.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	if err := w.WriteLine("hello stdout"); err != nil {
		t.Errorf("WriteLine failed: %v", err)
	}
}

func TestNew_WithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	w, err := output.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Errorf("WriteLine(%q) failed: %v", l, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	got := string(data)
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected file to contain %q, got:\n%s", l, got)
		}
	}
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := output.New("/nonexistent/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestWriteLine_Concurrent(t *testing.T) {
	w, err := output.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func(n int) {
			w.WriteLine("concurrent line")
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
