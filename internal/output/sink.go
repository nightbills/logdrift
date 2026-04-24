package output

import (
	"context"
)

// Sink drains a channel of formatted lines and writes each one via Writer
// until the channel is closed or the context is cancelled.
//
// It returns the first write error encountered, or nil on clean shutdown.
func Sink(ctx context.Context, w *Writer, lines <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-lines:
			if !ok {
				// Channel closed — all producers finished.
				return nil
			}
			if err := w.WriteLine(line); err != nil {
				return err
			}
		}
	}
}
