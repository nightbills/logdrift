// Package fieldmask implements field-level projection for structured JSON logs.
//
// When tailing logs from multiple services, it is often useful to focus on a
// subset of fields or to suppress noisy metadata. fieldmask supports two modes:
//
//   - ModeInclude: render only the explicitly listed fields (allowlist).
//   - ModeExclude: render all fields except the listed ones (denylist).
//
// Usage:
//
//	// Show only level and msg
//	m := fieldmask.New("level,msg", fieldmask.ModeInclude)
//	filtered := m.Apply(rawJSONLine)
//
//	// Hide the ts and trace_id fields
//	m := fieldmask.New("ts,trace_id", fieldmask.ModeExclude)
//	filtered := m.Apply(rawJSONLine)
//
// Non-JSON lines are passed through unchanged so that plain-text log entries
// from mixed-format services are never silently dropped.
package fieldmask
