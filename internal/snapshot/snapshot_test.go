package snapshot

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAdd_And_Entries_ReturnLines(t *testing.T) {
	s := New(0)
	s.Add("line one")
	s.Add("line two")
	entries := s.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Line != "line one" {
		t.Errorf("unexpected first line: %q", entries[0].Line)
	}
}

func TestPrune_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	s := New(5 * time.Second)
	s.now = fixedClock(now)
	s.Add("old line")

	s.now = fixedClock(now.Add(10 * time.Second))
	s.Add("new line")

	entries := s.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after pruning, got %d", len(entries))
	}
	if entries[0].Line != "new line" {
		t.Errorf("unexpected line after pruning: %q", entries[0].Line)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	s := New(0)
	s.Add("a")
	s.Add("b")
	s.Reset()
	if len(s.Entries()) != 0 {
		t.Error("expected empty snapshot after Reset")
	}
}

func TestZeroTTL_KeepsAllEntries(t *testing.T) {
	s := New(0)
	for i := 0; i < 10; i++ {
		s.Add(fmt.Sprintf("line %d", i))
	}
	if len(s.Entries()) != 10 {
		t.Errorf("expected 10 entries with zero TTL, got %d", len(s.Entries()))
	}
}

func TestAdd_Concurrent_DoesNotRace(t *testing.T) {
	s := New(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.Add(fmt.Sprintf("goroutine %d", n))
		}(i)
	}
	wg.Wait()
	if len(s.Entries()) == 0 {
		t.Error("expected entries after concurrent adds")
	}
}
