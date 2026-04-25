// Package dedupe provides a simple deduplication filter for log lines.
// It suppresses repeated identical log lines within a configurable window,
// helping reduce noise when a service emits the same message repeatedly.
package dedupe

import (
	"sync"
	"time"
)

// Filter tracks recently seen log lines and reports whether a line is a
// duplicate within the configured TTL window.
type Filter struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	ttl     time.Duration
	stopCh  chan struct{}
}

// New creates a Filter that suppresses duplicate lines seen within ttl.
// A background goroutine periodically evicts expired entries; call Stop
// when the filter is no longer needed.
func New(ttl time.Duration) *Filter {
	f := &Filter{
		seen:   make(map[string]time.Time),
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}
	go f.evictLoop()
	return f
}

// IsDuplicate returns true if line was already seen within the TTL window.
// If it is not a duplicate, the line is recorded and false is returned.
func (f *Filter) IsDuplicate(line string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	if exp, ok := f.seen[line]; ok && now.Before(exp) {
		return true
	}
	f.seen[line] = now.Add(f.ttl)
	return false
}

// Stop halts the background eviction goroutine.
func (f *Filter) Stop() {
	close(f.stopCh)
}

func (f *Filter) evictLoop() {
	ticker := time.NewTicker(f.ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f.evict()
		case <-f.stopCh:
			return
		}
	}
}

func (f *Filter) evict() {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for k, exp := range f.seen {
		if now.After(exp) {
			delete(f.seen, k)
		}
	}
}
