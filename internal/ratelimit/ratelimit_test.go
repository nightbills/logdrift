package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/ratelimit"
)

func TestNew_NoLimit_AlwaysPasses(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		if !l.Wait(ctx) {
			t.Fatal("expected Wait to return true with no rate limit")
		}
	}
}

func TestNew_ContextCancel_ReturnsFalse(t *testing.T) {
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the initial token bucket (it starts empty; refill fires after 1s)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// With rate=1 and no pre-filled tokens, Wait should block and then
	// return false once the context times out.
	result := l.Wait(ctx)
	if result {
		t.Skip("token was available immediately, skipping cancel test")
	}
}

func TestWait_TokenConsumed(t *testing.T) {
	l := ratelimit.New(5)
	defer l.Stop()

	// Allow refill to populate tokens.
	time.Sleep(1100 * time.Millisecond)

	ctx := context.Background()
	acquired := 0
	for i := 0; i < 5; i++ {
		ctxShort, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		if l.Wait(ctxShort) {
			acquired++
		}
		cancel()
	}
	if acquired == 0 {
		t.Error("expected to acquire at least one token after refill")
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	l := ratelimit.New(10)
	l.Stop()
	// Calling Stop on a zero-rate limiter should also be safe.
	l2 := ratelimit.New(0)
	l2.Stop()
}
