// Package levelfilter implements ordered log-level filtering for logdrift.
//
// Log levels are ordered from least to most severe:
//
//	debug < info < warn < error < fatal
//
// Usage:
//
//	f := levelfilter.New(levelfilter.LevelWarn)
//	if f.Allow(entry["level"]) {
//	    // forward the log line
//	}
//
// Lines whose level field is absent or unrecognised are always forwarded so
// that non-standard log producers are never silently suppressed.
package levelfilter
