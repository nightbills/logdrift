package linecount_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/user/logdrift/internal/linecount"
)

func TestInc_And_Get(t *testing.T) {
	c := linecount.New()
	c.Inc("api")
	c.Inc("api")
	c.Inc("worker")

	if got := c.Get("api"); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
	if got := c.Get("worker"); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGet_UnknownService_ReturnsZero(t *testing.T) {
	c := linecount.New()
	if got := c.Get("unknown"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestReset_ClearsAllCounts(t *testing.T) {
	c := linecount.New()
	c.Inc("api")
	c.Inc("db")
	c.Reset()

	if got := c.Get("api"); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
	if got := c.Get("db"); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}

func TestSummary_ContainsAllServices(t *testing.T) {
	c := linecount.New()
	c.Inc("alpha")
	c.Inc("beta")
	c.Inc("beta")

	summary := c.Summary()
	if !strings.Contains(summary, "alpha: 1") {
		t.Errorf("summary missing alpha: 1, got:\n%s", summary)
	}
	if !strings.Contains(summary, "beta: 2") {
		t.Errorf("summary missing beta: 2, got:\n%s", summary)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	c := linecount.New()
	c.Inc("svc")
	snap := c.Snapshot()
	c.Inc("svc") // should not affect snap

	if snap["svc"] != 1 {
		t.Errorf("expected snapshot value 1, got %d", snap["svc"])
	}
}

func TestInc_ConcurrentSafe(t *testing.T) {
	c := linecount.New()
	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("shared")
		}()
	}
	wg.Wait()

	if got := c.Get("shared"); got != goroutines {
		t.Errorf("expected %d, got %d", goroutines, got)
	}
}
