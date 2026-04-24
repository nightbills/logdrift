package filter

import (
	"encoding/json"
	"strings"
)

// Options holds the filtering criteria for log entries.
type Options struct {
	Level   string // e.g. "error", "warn", "info"
	Service string // service name to match
	Contains string // substring to search for in the raw log line
}

// Entry represents a parsed JSON log entry.
type Entry map[string]interface{}

// Parse attempts to decode a raw log line into an Entry.
// Returns nil if the line is not valid JSON.
func Parse(line string) Entry {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '{' {
		return nil
	}
	var entry Entry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil
	}
	return entry
}

// Match reports whether the log entry satisfies all non-empty filter options.
func Match(entry Entry, raw string, opts Options) bool {
	if opts.Contains != "" && !strings.Contains(raw, opts.Contains) {
		return false
	}

	if opts.Level != "" {
		if entry == nil {
			return false
		}
		level := stringField(entry, "level", "lvl", "severity")
		if !strings.EqualFold(level, opts.Level) {
			return false
		}
	}

	if opts.Service != "" {
		if entry == nil {
			return false
		}
		service := stringField(entry, "service", "app", "name")
		if !strings.EqualFold(service, opts.Service) {
			return false
		}
	}

	return true
}

// stringField returns the string value of the first matching key found in entry.
func stringField(entry Entry, keys ...string) string {
	for _, k := range keys {
		if v, ok := entry[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}
