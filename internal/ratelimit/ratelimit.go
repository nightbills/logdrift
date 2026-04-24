// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the number of log lines emitted per second. This is useful
// when tailing high-throughput services to prevent terminal flooding.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which log lines are passed through.
type Limiter struct {
	tokens   chan struct{}
	rate     int
	stop     chan struct{}
}

// New creates a Limiter that allows up to ratePerSec lines per second.
// If ratePerSec is 0 or negative, no rate limiting is applied.
func New(ratePerSec int) *Limiter {
	l := &Limiter{
		rate: ratePerSec,
		stop: make(chan struct{}),
	}
	if ratePerSec > 0 {
		l.tokens = make(chan struct{}, ratePerSec)
		go l.refill()
	}
	return l
}

// refill adds tokens to the bucket at the configured rate, once per second.
func (l *Limiter) refill() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			for i := 0; i < l.rate; i++ {
				select {
				case l.tokens <- struct{}{}:
				default:
				}
			}
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
// Returns false if the context was cancelled before a token was acquired.
func (l *Limiter) Wait(ctx context.Context) bool {
	if l.rate <= 0 {
		return true
	}
	select {
	case <-ctx.Done():
		return false
	case <-l.tokens:
		return true
	}
}

// Stop shuts down the background refill goroutine.
func (l *Limiter) Stop() {
	if l.rate > 0 {
		close(l.stop)
	}
}
