// Package multiline collapses multi-line log entries into a single logical
// line before they are passed to the rest of the logdrift pipeline.
//
// # Continuation Detection
//
// A line is considered a continuation of the previous entry when it begins
// with one or more ASCII space or tab characters. This matches the convention
// used by most log frameworks when emitting stack traces or wrapped JSON.
//
// # Flush Conditions
//
// A buffered entry is flushed (emitted on the Out channel) when any of the
// following conditions are met:
//
//   - A new non-continuation line arrives (the previous entry is complete).
//   - The number of accumulated lines reaches the configured maxLines limit.
//   - The maxWait timer fires (useful for entries that arrive with no
//     following line to trigger a natural flush).
//   - Flush is called explicitly.
//
// Set maxLines to 0 to disable the line-count limit.
// Set maxWait to 0 to disable the timer-based flush.
package multiline
