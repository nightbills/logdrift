// Package snapshot provides a short-term, time-bounded in-memory store for
// recently seen log lines. It is used to display a rolling "live snapshot"
// of activity across all tailed services, giving operators a quick view of
// what happened in the last N seconds without scrolling back through output.
//
// # Overview
//
// A Snapshot holds a fixed-size window of log entries, each tagged with the
// wall-clock time at which it was received. Entries older than the configured
// TTL are pruned automatically on every write and can also be pruned on demand
// via Prune. The entire store can be cleared with Reset.
//
// # Usage
//
//	s := snapshot.New(30*time.Second, time.Now)
//	s.Add("service-a", line)
//	entries := s.Entries()
//	for _, e := range entries {
//		fmt.Println(e.Service, e.Line)
//	}
//
// Snapshot is safe for concurrent use.
package snapshot
