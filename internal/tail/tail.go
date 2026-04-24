package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single log line read from a source.
type Line struct {
	Service string
	Text    string
	Err     error
}

// TailFile tails the given file, sending new lines to the returned channel.
// It follows the file (like `tail -f`) until the context is cancelled.
func TailFile(ctx context.Context, path, service string) (<-chan Line, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan Line, 64)

	go func() {
		defer close(ch)
		defer f.Close()

		// Seek to end so we only tail new lines.
		if _, err := f.Seek(0, io.SeekEnd); err != nil {
			ch <- Line{Service: service, Err: err}
			return
		}

		reader := bufio.NewReader(f)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// No new data yet; wait briefly and retry.
					time.Sleep(200 * time.Millisecond)
					continue
				}
				ch <- Line{Service: service, Err: err}
				return
			}

			if line != "" {
				ch <- Line{Service: service, Text: line}
			}
		}
	}()

	return ch, nil
}

// MergeLines fans-in multiple Line channels into a single channel.
func MergeLines(ctx context.Context, sources ...<-chan Line) <-chan Line {
	merged := make(chan Line, 128)

	for _, src := range sources {
		go func(ch <-chan Line) {
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-ch:
					if !ok {
						return
					}
					merged <- line
				}
			}
		}(src)
	}

	return merged
}
