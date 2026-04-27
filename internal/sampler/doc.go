// Package sampler implements per-key rate-based sampling for log lines.
//
// When logdrift is tailing a very high-volume service it can be useful to
// see only a representative fraction of lines rather than every entry.
// Sampler tracks an independent counter per service (or any string key) and
// returns true from Allow only on the 1st, (N+1)th, (2N+1)th, … call.
//
// Example usage:
//
//	s := sampler.New(10)          // emit every 10th line per service
//	for line := range lines {
//		if s.Allow(line.Service) {
//			fmt.Println(line.Text)
//		}
//	}
//
// A rate of 1 (the default when 0 is supplied) disables sampling so that
// every line is forwarded, making it safe to always construct a Sampler
// regardless of whether the user requested sampling.
package sampler
