package sampler_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/sampler"
)

func TestNew_RateOne_AllowsAll(t *testing.T) {
	s := sampler.New(1)
	for i := 0; i < 10; i++ {
		if !s.Allow("svc") {
			t.Fatalf("expected Allow to return true at iteration %d", i)
		}
	}
}

func TestNew_RateZero_TreatedAsOne(t *testing.T) {
	s := sampler.New(0)
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %d", s.Rate())
	}
	if !s.Allow("svc") {
		t.Fatal("expected Allow to return true")
	}
}

func TestAllow_RateN_EmitsEveryNth(t *testing.T) {
	s := sampler.New(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Allow("svc")
	}
	// calls 1,4,7 (indices 0,3,6) should be true
	expected := []bool{true, false, false, true, false, false, true, false, false}
	for i, want := range expected {
		if results[i] != want {
			t.Errorf("index %d: got %v, want %v", i, results[i], want)
		}
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	s := sampler.New(2)
	// first call for each key should be allowed
	if !s.Allow("alpha") {
		t.Error("expected alpha first call to be allowed")
	}
	if !s.Allow("beta") {
		t.Error("expected beta first call to be allowed")
	}
	// second calls should be suppressed
	if s.Allow("alpha") {
		t.Error("expected alpha second call to be suppressed")
	}
	if s.Allow("beta") {
		t.Error("expected beta second call to be suppressed")
	}
}

func TestReset_RestartsCounters(t *testing.T) {
	s := sampler.New(3)
	s.Allow("svc") // count=1 → true
	s.Allow("svc") // count=2 → false
	s.Reset()
	if !s.Allow("svc") {
		t.Error("expected Allow to return true after Reset")
	}
}

func TestAllow_Concurrent_DoesNotPanic(t *testing.T) {
	s := sampler.New(5)
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			s.Allow("concurrent")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
