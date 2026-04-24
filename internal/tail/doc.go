// Package tail provides utilities for following log files in real time
// and merging line streams from multiple services into a single channel.
//
// Usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	ch1, err := tail.TailFile(ctx, "/var/log/api.log", "api")
//	ch2, err := tail.TailFile(ctx, "/var/log/worker.log", "worker")
//
//	for line := range tail.MergeLines(ctx, ch1, ch2) {
//		if line.Err != nil {
//			log.Println("tail error:", line.Err)
//			continue
//		}
//		fmt.Printf("[%s] %s", line.Service, line.Text)
//	}
//
// TailFile seeks to the end of the file before reading, so only lines
// appended after the call are emitted. It polls for new content every
// 200 ms when EOF is reached.
package tail
