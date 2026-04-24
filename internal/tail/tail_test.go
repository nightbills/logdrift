package tail

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logdrift-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestTailFile_ReceivesNewLines(t *testing.T) {
	path := writeTempLog(t, "")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch, err := TailFile(ctx, path, "svc-a")
	if err != nil {
		t.Fatalf("TailFile returned error: %v", err)
	}

	// Append a line after tailing starts.
	go func() {
		time.Sleep(100 * time.Millisecond)
		f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(`{"level":"info","msg":"hello"}` + "\n")
	}()

	select {
	case line := <-ch:
		if line.Err != nil {
			t.Fatalf("unexpected error: %v", line.Err)
		}
		if line.Service != "svc-a" {
			t.Errorf("expected service svc-a, got %q", line.Service)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for line")
	}
}

func TestTailFile_InvalidPath(t *testing.T) {
	ctx := context.Background()
	_, err := TailFile(ctx, "/nonexistent/path/to/file.log", "svc")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestMergeLines_FansIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := make(chan Line, 1)
	ch2 := make(chan Line, 1)
	ch1 <- Line{Service: "a", Text: "line-a"}
	ch2 <- Line{Service: "b", Text: "line-b"}

	merged := MergeLines(ctx, ch1, ch2)

	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case l := <-merged:
			seen[l.Service] = true
		case <-ctx.Done():
			t.Fatal("timed out waiting for merged lines")
		}
	}

	if !seen["a"] || !seen["b"] {
		t.Errorf("expected both services, got %v", seen)
	}
}
