package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/user/logdrift/internal/colorscheme"
	"github.com/user/logdrift/internal/dedupe"
	"github.com/user/logdrift/internal/fieldmask"
	"github.com/user/logdrift/internal/filter"
	"github.com/user/logdrift/internal/formatter"
	"github.com/user/logdrift/internal/highlight"
	"github.com/user/logdrift/internal/levelfilter"
	"github.com/user/logdrift/internal/linecount"
	"github.com/user/logdrift/internal/multiline"
	"github.com/user/logdrift/internal/output"
	"github.com/user/logdrift/internal/ratelimit"
	"github.com/user/logdrift/internal/redact"
	"github.com/user/logdrift/internal/retag"
	"github.com/user/logdrift/internal/tail"
	"github.com/user/logdrift/internal/timerange"
	"github.com/user/logdrift/internal/truncate"
)

func main() {
	flags := parseFlags()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	out, err := output.New(flags.outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logdrift: output: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	scheme := colorscheme.Get(flags.colorScheme)
	filters, err := filter.Parse(flags.filterExpr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logdrift: filter: %v\n", err)
		os.Exit(1)
	}

	lf := levelfilter.New(levelfilter.Parse(flags.minLevel))
	tr, _ := timerange.New(flags.fromTime, flags.toTime)
	rl := ratelimit.New(flags.rateLimit)
	defer rl.Stop()
	dd := dedupe.New(5 * time.Second)
	defer dd.Stop()
	rd := redact.New(splitCSV(flags.redactFields), "[REDACTED]")
	fm := fieldmask.New(splitCSV(flags.includeFields), splitCSV(flags.excludeFields))
	rt, _ := retag.New(nil)
	hl := highlight.New(splitCSV(flags.keywords), flags.caseSensitive)
	trunc := truncate.New(flags.maxWidth)
	lc := linecount.New()
	ml := multiline.New(flags.multilineMax, time.Duration(flags.multilineWaitMS)*time.Millisecond, " ")

	var channels []<-chan string
	for _, path := range flags.files {
		ch, err := tail.TailFile(ctx, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logdrift: tail %s: %v\n", path, err)
			continue
		}
		channels = append(channels, ch)
	}
	if len(channels) == 0 {
		fmt.Fprintln(os.Stderr, "logdrift: no files to tail")
		os.Exit(1)
	}

	merged := tail.MergeLines(ctx, channels...)

	go func() {
		for {
			select {
			case line, ok := <-merged:
				if !ok {
					ml.Flush()
					return
				}
				ml.Feed(line)
			case <-ctx.Done():
				ml.Flush()
				return
			}
		}
	}()

	for {
		select {
		case line, ok := <-ml.Out():
			if !ok {
				return
			}
			if !filter.Match(line, filters) {
				continue
			}
			if !lf.Allow(line) {
				continue
			}
			if tr != nil && !tr.Allow(line) {
				continue
			}
			if dd.IsDuplicate(line) {
				continue
			}
			if !rl.Wait(ctx) {
				continue
			}
			line = rd.Apply(line)
			line = fm.Apply(line)
			line = rt.Apply(line)
			line = formatter.Format(line, scheme)
			line = hl.Apply(line)
			line = trunc.Line(line)
			lc.Inc("default")
			out.WriteLine(line)
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "\nlogdrift: %s\n", lc.Summary())
			return
		}
	}
}

type flags struct {
	files           []string
	filterExpr      string
	minLevel        string
	outputFile      string
	colorScheme     string
	keywords        string
	caseSensitive   bool
	maxWidth        int
	rateLimit       int
	redactFields    string
	includeFields   string
	excludeFields   string
	fromTime        string
	toTime          string
	multilineMax    int
	multilineWaitMS int
}

func parseFlags() flags {
	var f flags
	flag.StringVar(&f.filterExpr, "filter", "", "filter expression (key=value)")
	flag.StringVar(&f.minLevel, "level", "debug", "minimum log level")
	flag.StringVar(&f.outputFile, "out", "", "write output to file")
	flag.StringVar(&f.colorScheme, "scheme", "default", "color scheme")
	flag.StringVar(&f.keywords, "highlight", "", "comma-separated keywords to highlight")
	flag.BoolVar(&f.caseSensitive, "case-sensitive", false, "case-sensitive keyword matching")
	flag.IntVar(&f.maxWidth, "width", 0, "truncate lines to width (0=off)")
	flag.IntVar(&f.rateLimit, "rate", 0, "max lines/sec (0=unlimited)")
	flag.StringVar(&f.redactFields, "redact", "", "comma-separated JSON fields to redact")
	flag.StringVar(&f.includeFields, "include-fields", "", "comma-separated fields to include")
	flag.StringVar(&f.excludeFields, "exclude-fields", "", "comma-separated fields to exclude")
	flag.StringVar(&f.fromTime, "from", "", "start time (RFC3339)")
	flag.StringVar(&f.toTime, "to", "", "end time (RFC3339)")
	flag.IntVar(&f.multilineMax, "multiline-max", 50, "max lines per multiline entry (0=unlimited)")
	flag.IntVar(&f.multilineWaitMS, "multiline-wait", 100, "max wait ms before flushing partial entry")
	flag.Parse()
	f.files = flag.Args()
	return f
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
