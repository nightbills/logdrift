// Package linecount tracks the number of log lines processed per service,
// providing a lightweight counter map safe for concurrent use.
package linecount

import (
	"fmt"
	"sort"
	"sync"
)

// Counter holds per-service line counts.
type Counter struct {
	mu     sync.Mutex
	counts map[string]int64
}

// New creates a new Counter.
func New() *Counter {
	return &Counter{
		counts: make(map[string]int64),
	}
}

// Inc increments the count for the given service by 1.
func (c *Counter) Inc(service string) {
	c.mu.Lock()
	c.counts[service]++
	c.mu.Unlock()
}

// Get returns the current count for the given service.
func (c *Counter) Get(service string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[service]
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.mu.Lock()
	c.counts = make(map[string]int64)
	c.mu.Unlock()
}

// Summary returns a formatted multi-line string listing each service
// and its line count, sorted alphabetically by service name.
func (c *Counter) Summary() string {
	c.mu.Lock()
	services := make([]string, 0, len(c.counts))
	for s := range c.counts {
		services = append(services, s)
	}
	c.mu.Unlock()

	sort.Strings(services)

	result := ""
	for _, s := range services {
		c.mu.Lock()
		n := c.counts[s]
		c.mu.Unlock()
		result += fmt.Sprintf("%s: %d\n", s, n)
	}
	return result
}

// Snapshot returns a copy of the current counts map.
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]int64, len(c.counts))
	for k, v := range c.counts {
		copy[k] = v
	}
	return copy
}
