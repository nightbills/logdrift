// Package multiline provides support for collapsing multi-line log entries
// into a single line before further processing. This is useful when log
// producers emit stack traces or JSON objects split across multiple lines.
package multiline

import (
	"strings"
	"time"
)

// Collapser accumulates lines that belong to a single logical log entry and
// emits them as one joined string once a flush condition is met.
type Collapser struct {
	maxLines  int
	maxWait   time.Duration
	delimiter string
	buf       []string
	timer     *time.Timer
	out       chan string
}

// New creates a Collapser. maxLines is the maximum number of raw lines to
// accumulate before forcing a flush (0 means unlimited). maxWait is the
// maximum time to hold a partial entry (0 means flush immediately on each
// line that is a continuation). delimiter is placed between joined lines.
func New(maxLines int, maxWait time.Duration, delimiter string) *Collapser {
	if delimiter == "" {
		delimiter = " "
	}
	return &Collapser{
		maxLines:  maxLines,
		maxWait:   maxWait,
		delimiter: delimiter,
		out:       make(chan string, 64),
	}
}

// Out returns the channel on which flushed entries are emitted.
func (c *Collapser) Out() <-chan string { return c.out }

// Feed accepts a raw line. Lines that start with whitespace are treated as
// continuations of the previous entry. All other lines trigger a flush of
// any buffered content before starting a new entry.
func (c *Collapser) Feed(line string) {
	isContinuation := len(line) > 0 && (line[0] == ' ' || line[0] == '\t')

	if !isContinuation && len(c.buf) > 0 {
		c.flush()
	}

	c.buf = append(c.buf, strings.TrimRight(line, "\r\n"))

	if c.maxLines > 0 && len(c.buf) >= c.maxLines {
		c.flush()
		return
	}

	if !isContinuation {
		if c.maxWait == 0 {
			c.flush()
			return
		}
		if c.timer != nil {
			c.timer.Reset(c.maxWait)
		} else {
			c.timer = time.AfterFunc(c.maxWait, func() { c.flush() })
		}
	}
}

// Flush forces any buffered lines to be emitted immediately.
func (c *Collapser) Flush() { c.flush() }

func (c *Collapser) flush() {
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	if len(c.buf) == 0 {
		return
	}
	c.out <- strings.Join(c.buf, c.delimiter)
	c.buf = c.buf[:0]
}
