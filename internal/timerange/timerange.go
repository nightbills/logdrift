// Package timerange provides filtering of log lines by a time window.
// It parses a "time" or "ts" field from JSON log entries and checks
// whether the timestamp falls within a [From, To] range.
package timerange

import (
	"encoding/json"
	"time"
)

// Filter holds an optional time window. A zero From or To means unbounded
// on that side.
type Filter struct {
	From time.Time
	To   time.Time
}

// New creates a Filter. Pass zero values to leave a bound open.
func New(from, to time.Time) *Filter {
	return &Filter{From: from, To: to}
}

// Allow returns true when the log line's timestamp is within the filter
// window, or when the line has no parseable timestamp.
// Recognised JSON fields: "time", "ts", "timestamp".
func (f *Filter) Allow(line string) bool {
	if f.From.IsZero() && f.To.IsZero() {
		return true
	}

	t, ok := extractTime(line)
	if !ok {
		// No parseable timestamp — let the line through.
		return true
	}

	if !f.From.IsZero() && t.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && t.After(f.To) {
		return false
	}
	return true
}

// extractTime attempts to read a timestamp from well-known JSON fields.
func extractTime(line string) (time.Time, bool) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return time.Time{}, false
	}

	for _, key := range []string{"time", "ts", "timestamp"} {
		v, ok := raw[key]
		if !ok {
			continue
		}

		// Try numeric Unix seconds / milliseconds first.
		var f float64
		if err := json.Unmarshal(v, &f); err == nil {
			if f > 1e12 { // milliseconds
				return time.UnixMilli(int64(f)).UTC(), true
			}
			return time.Unix(int64(f), 0).UTC(), true
		}

		// Try RFC3339 string.
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
				if t, err := time.Parse(layout, s); err == nil {
					return t.UTC(), true
				}
			}
		}
	}
	return time.Time{}, false
}
