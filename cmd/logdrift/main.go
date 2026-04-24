// main is the entry point for the logdrift CLI tool.
// It wires together the tail, filter, and formatter packages
// to provide real-time structured JSON log viewing.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logdrift/internal/filter"
	"github.com/yourorg/logdrift/internal/formatter"
	"github.com/yourorg/logdrift/internal/tail"
)

// config holds the parsed CLI flags.
type config struct {
	files   []string
	level   string
	service string
	fields  map[string]string
	noColor bool
}

func main() {
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.files) == 0 {
		fmt.Fprintln(os.Stderr, "error: at least one log file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Build filter expressions from flags.
	var exprs []filter.Expression
	if cfg.level != "" {
		exprs = append(exprs, filter.Expression{Field: "level", Value: cfg.level})
	}
	if cfg.service != "" {
		exprs = append(exprs, filter.Expression{Field: "service", Value: cfg.service})
	}
	for k, v := range cfg.fields {
		exprs = append(exprs, filter.Expression{Field: k, Value: v})
	}

	// Start tailing each file and merge into a single channel.
	channels := make([]<-chan string, 0, len(cfg.files))
	for _, path := range cfg.files {
		ch, err := tail.TailFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error tailing %s: %v\n", path, err)
			os.Exit(1)
		}
		channels = append(channels, ch)
	}

	merged := tail.MergeLines(channels...)

	fmtOpts := formatter.Options{
		NoColor: cfg.noColor,
	}

	for line := range merged {
		entry, err := filter.Parse(line)
		if err != nil {
			// Non-JSON lines are passed through as-is.
			fmt.Println(line)
			continue
		}
		if !filter.Match(entry, exprs) {
			continue
		}
		fmt.Println(formatter.Format(line, fmtOpts))
	}
}

// parseFlags parses command-line arguments and returns a config.
func parseFlags(args []string) (config, error) {
	fs := flag.NewFlagSet("logdrift", flag.ContinueOnError)

	level := fs.String("level", "", "filter by log level (e.g. error, warn, info)")
	service := fs.String("service", "", "filter by service name")
	rawFields := fs.String("fields", "", "additional field filters as key=value pairs, comma-separated")
	noColor := fs.Bool("no-color", false, "disable colored output")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: logdrift [flags] <file> [file...]")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	fields := make(map[string]string)
	if *rawFields != "" {
		for _, pair := range strings.Split(*rawFields, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(parts) != 2 {
				return config{}, fmt.Errorf("invalid field filter %q: expected key=value", pair)
			}
			fields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return config{
		files:   fs.Args(),
		level:   strings.ToLower(*level),
		service: *service,
		fields:  fields,
		noColor: *noColor,
	}, nil
}
