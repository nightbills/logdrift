package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/throttle"
)

func TestAllow_ZeroRate_AlwaysPasses(t *testing.T) {
	th := throttle.New(0)
	defer th.Stop()
	for i := 0; i < 1000; i++ {
		if !th.Allow("svc") {
			t.Fatal("expected Allow to return true for zero rate")
		}
	}
}

func TestAllow_RateOne_FirstPasses(t *testing.T) {
	th := throttle.New(1)
	defer th.Stop()
	if !th.Allow("svc") {
		t.Fatal("first call should pass")
	}
	if th.Allow("svc") {
		t.Fatal("second call within same second should be dropped")
	}
}

func TestAllow_RateN_AllowsExactlyN(t *testing.T) {
	const rate = 5
	th := throttle.New(rate)
	defer th.Stop()
	passed := 0
	for i := 0; i < 10; i++ {
		if th.Allow("svc") {
			passed++
		}
	}
	if passed != rate {
		t.Fatalf("expected %d allowed, got %d", rate, passed)
	}
}

func TestDropped_CountsExcess(t *testing.T) {
	th := throttle.New(2)
	defer th.Stop()
	for i := 0; i < 5; i++ {
		th.Allow("svc")
	}
	if got := th.Dropped("svc"); got != 3 {
		t.Fatalf("expected 3 dropped, got %d", got)
	}
}

func TestDropped_UnknownKey_ReturnsZero(t *testing.T) {
	th := throttle.New(10)
	defer th.Stop()
	if d := th.Dropped("unknown"); d != 0 {
		t.Fatalf("expected 0, got %d", d)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(1)
	defer th.Stop()
	if !th.Allow("a") {
		t.Fatal("key a first call should pass")
	}
	if !th.Allow("b") {
		t.Fatal("key b first call should pass independently")
	}
	if th.Allow("a") {
		t.Fatal("key a second call should be dropped")
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	th := throttle.New(10)
	th.Allow("svc")
	th.Stop()
	time.Sleep(10 * time.Millisecond) // ensure goroutine exits cleanly
}
