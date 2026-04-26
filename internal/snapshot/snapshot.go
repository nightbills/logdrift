// Package snapshot provides periodic capture and replay of log lines
// seen within a rolling time window, useful for summarising bursts of
// activity without re-tailing files from scratch.
package snapshot

import (
	"sync"
	"time"
)

// Entry holds a single captured log line together with the time it was
// received by logdrift.
type Entry struct {
	Line      string
	ReceivedAt time.Time
}

// Snapshot accumulates log lines in memory and discards entries older
// than the configured TTL.
type Snapshot struct {
	mu      sync.Mutex
	entries []Entry
	ttl     time.Duration
	now     func() time.Time // injectable for testing
}

// New returns a Snapshot that retains entries for the given TTL.
// A zero TTL means entries are kept indefinitely.
func New(ttl time.Duration) *Snapshot {
	return &Snapshot{
		ttl: ttl,
		now: time.Now,
	}
}

// Add appends a line to the snapshot, pruning stale entries first.
func (s *Snapshot) Add(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune()
	s.entries = append(s.entries, Entry{Line: line, ReceivedAt: s.now()})
}

// Entries returns a copy of the current in-window entries.
func (s *Snapshot) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Reset clears all retained entries.
func (s *Snapshot) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}

// prune removes entries older than the TTL. Caller must hold s.mu.
func (s *Snapshot) prune() {
	if s.ttl == 0 {
		return
	}
	cutoff := s.now().Add(-s.ttl)
	i := 0
	for ; i < len(s.entries); i++ {
		if s.entries[i].ReceivedAt.After(cutoff) {
			break
		}
	}
	s.entries = s.entries[i:]
}
