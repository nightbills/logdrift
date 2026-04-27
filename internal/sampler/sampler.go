// Package sampler provides log-line sampling to reduce high-volume output
// by emitting only every Nth matching line per service.
package sampler

import (
	"sync"
)

// Sampler tracks per-service counters and emits a line only when the counter
// is a multiple of the configured rate. A rate of 1 (or 0) passes every line.
type Sampler struct {
	mu      sync.Mutex
	rate    int
	counts  map[string]int
}

// New creates a Sampler with the given sample rate. A rate <= 1 disables
// sampling and every line is allowed through.
func New(rate int) *Sampler {
	if rate < 1 {
		rate = 1
	}
	return &Sampler{
		rate:   rate,
		counts: make(map[string]int),
	}
}

// Allow increments the counter for the given key (typically a service name)
// and returns true when the line should be emitted.
// With rate=1 every call returns true.
// With rate=N the first call and every Nth subsequent call returns true.
func (s *Sampler) Allow(key string) bool {
	if s.rate == 1 {
		return true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key]++
	return s.counts[key]%s.rate == 1
}

// Reset clears all counters, restarting sampling from zero for every key.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() int {
	return s.rate
}
