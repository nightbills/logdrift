// Package levelfilter provides ordered log level filtering so that
// only messages at or above a minimum severity are forwarded.
package levelfilter

import "strings"

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

// Parse converts a level string (case-insensitive) to a Level.
// Returns LevelUnknown if the string is not recognised.
func Parse(s string) Level {
	if l, ok := levelNames[strings.ToLower(s)]; ok {
		return l
	}
	return LevelUnknown
}

// String returns the canonical name for a Level.
func (l Level) String() string {
	for name, lvl := range levelNames {
		if lvl == l {
			return name
		}
	}
	return "unknown"
}

// Filter holds the minimum level threshold.
type Filter struct {
	min Level
}

// New creates a Filter that passes log lines whose level is >= minLevel.
// If minLevel is LevelUnknown the filter passes everything.
func New(minLevel Level) *Filter {
	return &Filter{min: minLevel}
}

// Allow returns true when the given level string meets the minimum threshold.
// Lines whose level cannot be parsed are always allowed through so that
// non-standard or missing level fields are never silently dropped.
func (f *Filter) Allow(levelStr string) bool {
	if f.min == LevelUnknown {
		return true
	}
	l := Parse(levelStr)
	if l == LevelUnknown {
		return true
	}
	return l >= f.min
}
