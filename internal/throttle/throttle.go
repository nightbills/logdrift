// Package throttle provides per-service line emission throttling.
// It limits how many lines per second are forwarded for a given service key,
// dropping excess lines and optionally reporting drop counts.
package throttle

import (
	"sync"
	"time"
)

// Throttler tracks per-key token buckets and drops lines that exceed the
// configured rate (lines per second). A rate of 0 disables throttling.
type Throttler struct {
	rate     int
	mu       sync.Mutex
	buckets  map[string]*bucket
	stopCh   chan struct{}
}

type bucket struct {
	tokens   int
	dropped  int
	lastTick time.Time
}

// New creates a Throttler that allows at most ratePerSec lines per second
// per key. Call Stop when done to release background resources.
func New(ratePerSec int) *Throttler {
	t := &Throttler{
		rate:    ratePerSec,
		buckets: make(map[string]*bucket),
		stopCh:  make(chan struct{}),
	}
	go t.refill()
	return t
}

// Allow returns true if the line for the given key should be forwarded.
// When rate is 0, Allow always returns true.
func (t *Throttler) Allow(key string) bool {
	if t.rate == 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	b, ok := t.buckets[key]
	if !ok {
		b = &bucket{tokens: t.rate, lastTick: time.Now()}
		t.buckets[key] = b
	}
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	b.dropped++
	return false
}

// Dropped returns the number of lines dropped for key since the last reset.
func (t *Throttler) Dropped(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if b, ok := t.buckets[key]; ok {
		return b.dropped
	}
	return 0
}

// Stop halts the background refill goroutine.
func (t *Throttler) Stop() {
	close(t.stopCh)
}

func (t *Throttler) refill() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.mu.Lock()
			for _, b := range t.buckets {
				b.tokens = t.rate
				b.dropped = 0
			}
			t.mu.Unlock()
		case <-t.stopCh:
			return
		}
	}
}
