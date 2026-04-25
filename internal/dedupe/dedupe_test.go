package dedupe_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/dedupe"
)

func TestIsDuplicate_FirstOccurrence_NotDuplicate(t *testing.T) {
	f := dedupe.New(1 * time.Second)
	defer f.Stop()

	if f.IsDuplicate("hello world") {
		t.Fatal("expected first occurrence to not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence_IsDuplicate(t *testing.T) {
	f := dedupe.New(1 * time.Second)
	defer f.Stop()

	f.IsDuplicate("repeated line")
	if !f.IsDuplicate("repeated line") {
		t.Fatal("expected second occurrence to be a duplicate")
	}
}

func TestIsDuplicate_DifferentLines_NotDuplicate(t *testing.T) {
	f := dedupe.New(1 * time.Second)
	defer f.Stop()

	f.IsDuplicate("line one")
	if f.IsDuplicate("line two") {
		t.Fatal("different lines should not be considered duplicates")
	}
}

func TestIsDuplicate_AfterTTLExpiry_NotDuplicate(t *testing.T) {
	ttl := 50 * time.Millisecond
	f := dedupe.New(ttl)
	defer f.Stop()

	f.IsDuplicate("expiring line")
	time.Sleep(ttl + 20*time.Millisecond)

	if f.IsDuplicate("expiring line") {
		t.Fatal("expected line to no longer be a duplicate after TTL expiry")
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	f := dedupe.New(100 * time.Millisecond)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Stop panicked: %v", r)
		}
	}()
	f.Stop()
}

func TestIsDuplicate_Concurrent(t *testing.T) {
	f := dedupe.New(500 * time.Millisecond)
	defer f.Stop()

	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			f.IsDuplicate("concurrent line")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
