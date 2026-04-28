package multiline_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/multiline"
)

// TestPipeline_StackTrace simulates a Java-style stack trace arriving as
// several indented continuation lines after the initial error message.
func TestPipeline_StackTrace(t *testing.T) {
	c := multiline.New(0, 10*time.Millisecond, " ")

	lines := []string{
		`{"level":"error","msg":"unhandled exception"}`,
		`\tat com.example.App.main(App.java:42)`,
		`\tat com.example.App.run(App.java:17)`,
		`{"level":"info","msg":"recovered"}`,
	}

	for _, l := range lines {
		c.Feed(l)
	}
	c.Flush()

	var got []string
	timeout := time.After(100 * time.Millisecond)
collect:
	for {
		select {
		case v := <-c.Out():
			got = append(got, v)
			if len(got) == 2 {
				break collect
			}
		case <-timeout:
			break collect
		}
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(got), got)
	}

	expected0 := `{"level":"error","msg":"unhandled exception"} \tat com.example.App.main(App.java:42) \tat com.example.App.run(App.java:17)`
	if got[0] != expected0 {
		t.Errorf("entry 0 mismatch:\n got:  %q\n want: %q", got[0], expected0)
	}

	if got[1] != `{"level":"info","msg":"recovered"}` {
		t.Errorf("entry 1 mismatch: %q", got[1])
	}
}

// TestPipeline_NoMultiline ensures plain single-line entries pass through
// unchanged and in order.
func TestPipeline_NoMultiline(t *testing.T) {
	c := multiline.New(0, 0, " ")
	input := []string{"alpha", "beta", "gamma"}
	for _, l := range input {
		c.Feed(l)
	}

	var got []string
	for i := 0; i < len(input); i++ {
		select {
		case v := <-c.Out():
			got = append(got, v)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timed out waiting for output")
		}
	}

	for i, want := range input {
		if got[i] != want {
			t.Errorf("line %d: got %q want %q", i, got[i], want)
		}
	}
}
