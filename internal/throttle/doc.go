// Package throttle implements per-key line-rate throttling for logdrift.
//
// A Throttler enforces a maximum number of log lines forwarded per second
// for each named key (typically a service name). Lines that exceed the
// configured rate are dropped, and the drop count is tracked so callers
// can surface a "N lines dropped" warning to the user.
//
// Usage:
//
//	th := throttle.New(100) // allow up to 100 lines/sec per service
//	defer th.Stop()
//
//	if th.Allow(serviceName) {
//		// forward the line
//	}
//
// A rate of 0 disables throttling entirely — all calls to Allow return true.
package throttle
