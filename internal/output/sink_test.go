package output_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/output"
)

func TestSink_WritesAllLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sink.log")

	w, err := output.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ch := make(chan string, 5)
	want := []string{"alpha", "beta", "gamma"}
	for _, l := range want {
		ch <- l
	}
	close(ch)

	ctx := context.Background()
	if err := output.Sink(ctx, w, ch); err != nil {
		t.Fatalf("Sink: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, _ := os.ReadFile(path)
	got := string(data)
	for _, l := range want {
		if !strings.Contains(got, l) {
			t.Errorf("missing %q in output:\n%s", l, got)
		}
	}
}

func TestSink_RespectsContextCancel(t *testing.T) {
	w, err := output.New("")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	ch := make(chan string) // never closed, never sent to

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	if err := output.Sink(ctx, w, ch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Error("Sink did not respect context cancellation in time")
	}
}
