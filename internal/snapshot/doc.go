// Package snapshot implements a rolling in-memory window of log lines
// received by logdrift during a live tail session.
//
// # Overview
//
// A [Snapshot] accumulates [Entry] values as lines arrive from one or more
// tailed files. Entries older than the configured TTL are lazily discarded
// whenever [Snapshot.Add] or [Snapshot.Entries] is called, keeping memory
// usage proportional to recent activity rather than total uptime.
//
// # Usage
//
//	snap := snapshot.New(30 * time.Second)
//
//	// Feed lines from the merged tail channel:
//	for line := range lines {
//		snap.Add(line)
//	}
//
//	// Later, replay the window (e.g. on SIGUSR1):
//	for _, e := range snap.Entries() {
//		fmt.Println(e.Line)
//	}
//
// A zero TTL retains all entries for the lifetime of the process.
package snapshot
