package timerange_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/logdrift/internal/timerange"
)

var (
	base  = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	early = base.Add(-time.Hour)
	late  = base.Add(time.Hour)
)

func rfc(t time.Time) string {
	return fmt.Sprintf(`{"time":%q,"msg":"hello"}`, t.Format(time.RFC3339))
}

func unix(t time.Time) string {
	return fmt.Sprintf(`{"ts":%d,"msg":"hello"}`, t.Unix())
}

func TestAllow_NoBounds_AlwaysTrue(t *testing.T) {
	f := timerange.New(time.Time{}, time.Time{})
	if !f.Allow(rfc(early)) {
		t.Fatal("expected true with no bounds")
	}
}

func TestAllow_WithinWindow_True(t *testing.T) {
	f := timerange.New(early, late)
	if !f.Allow(rfc(base)) {
		t.Fatal("expected true for timestamp inside window")
	}
}

func TestAllow_BeforeFrom_False(t *testing.T) {
	f := timerange.New(base, late)
	if f.Allow(rfc(early)) {
		t.Fatal("expected false for timestamp before From")
	}
}

func TestAllow_AfterTo_False(t *testing.T) {
	f := timerange.New(early, base)
	if f.Allow(rfc(late)) {
		t.Fatal("expected false for timestamp after To")
	}
}

func TestAllow_UnixTimestamp_Parsed(t *testing.T) {
	f := timerange.New(early, late)
	if !f.Allow(unix(base)) {
		t.Fatal("expected true for unix timestamp inside window")
	}
}

func TestAllow_NonJSON_PassesThrough(t *testing.T) {
	f := timerange.New(early, late)
	if !f.Allow("not json at all") {
		t.Fatal("expected non-JSON lines to pass through")
	}
}

func TestAllow_NoTimestampField_PassesThrough(t *testing.T) {
	f := timerange.New(early, late)
	if !f.Allow(`{"msg":"no time field here"}`) {
		t.Fatal("expected lines without timestamp to pass through")
	}
}

func TestAllow_OnlyFromBound(t *testing.T) {
	f := timerange.New(base, time.Time{})
	if f.Allow(rfc(early)) {
		t.Fatal("expected false: timestamp before From-only bound")
	}
	if !f.Allow(rfc(late)) {
		t.Fatal("expected true: timestamp after From-only bound")
	}
}

func TestAllow_OnlyToBound(t *testing.T) {
	f := timerange.New(time.Time{}, base)
	if !f.Allow(rfc(early)) {
		t.Fatal("expected true: timestamp before To-only bound")
	}
	if f.Allow(rfc(late)) {
		t.Fatal("expected false: timestamp after To-only bound")
	}
}
