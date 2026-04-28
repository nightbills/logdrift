package multiline

import (
	"testing"
	"time"
)

func drain(ch <-chan string) []string {
	var out []string
	for {
		select {
		case v := <-ch:
			out = append(out, v)
		default:
			return out
		}
	}
}

func TestFeed_SingleLine_EmittedImmediately(t *testing.T) {
	c := New(0, 0, " ")
	c.Feed("hello world")
	got := drain(c.Out())
	if len(got) != 1 || got[0] != "hello world" {
		t.Fatalf("expected [hello world], got %v", got)
	}
}

func TestFeed_ContinuationLine_Collapsed(t *testing.T) {
	c := New(0, 0, " ")
	c.Feed("first line")
	c.Feed("  continuation")
	c.Feed("next entry")
	got := drain(c.Out())
	if len(got) != 1 {
		t.Fatalf("expected 1 flushed entry before next, got %d: %v", len(got), got)
	}
	if got[0] != "first line   continuation" {
		t.Errorf("unexpected collapsed line: %q", got[0])
	}
}

func TestFeed_MaxLines_ForcesFlush(t *testing.T) {
	c := New(2, 0, "|")
	c.Feed("line1")
	c.Feed("  line2") // continuation but maxLines reached
	got := drain(c.Out())
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0] != "line1|  line2" {
		t.Errorf("unexpected content: %q", got[0])
	}
}

func TestFeed_MaxWait_FlushesAfterTimeout(t *testing.T) {
	c := New(0, 20*time.Millisecond, " ")
	c.Feed("standalone")
	time.Sleep(50 * time.Millisecond)
	got := drain(c.Out())
	if len(got) != 1 || got[0] != "standalone" {
		t.Fatalf("expected [standalone], got %v", got)
	}
}

func TestFlush_EmptyBuffer_NoOutput(t *testing.T) {
	c := New(0, 0, " ")
	c.Flush()
	got := drain(c.Out())
	if len(got) != 0 {
		t.Errorf("expected no output, got %v", got)
	}
}

func TestFeed_CustomDelimiter(t *testing.T) {
	c := New(0, 0, "\\n")
	c.Feed("msg")
	c.Feed("  trace line")
	c.Flush()
	got := drain(c.Out())
	if len(got) != 1 || got[0] != "msg\\n  trace line" {
		t.Errorf("unexpected: %v", got)
	}
}
