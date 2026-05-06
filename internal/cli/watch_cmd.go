package cli

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"envdiff/internal/diff"
)

// WatchArgs holds parsed arguments for the watch sub-command.
type WatchArgs struct {
	LeftPath  string
	RightPath string
	Interval  time.Duration
	Format    diff.Format
}

// parseWatchArgs extracts watch sub-command arguments from os.Args style slice.
// Expected: watch [--interval=Ns] [--format=fmt] <left> <right>
func parseWatchArgs(args []string) (WatchArgs, error) {
	wa := WatchArgs{
		Interval: 2 * time.Second,
		Format:   diff.FormatText,
	}
	positional := []string{}
	for _, a := range args {
		var val string
		if matchFlag(a, "--interval", &val) {
			d, err := time.ParseDuration(val)
			if err != nil {
				return wa, fmt.Errorf("invalid interval %q: %w", val, err)
			}
			wa.Interval = d
		} else if matchFlag(a, "--format", &val) {
			f, err := diff.ParseFormat(val)
			if err != nil {
				return wa, err
			}
			wa.Format = f
		} else {
			positional = append(positional, a)
		}
	}
	if len(positional) < 2 {
		return wa, fmt.Errorf("watch requires two file paths")
	}
	wa.LeftPath = positional[0]
	wa.RightPath = positional[1]
	return wa, nil
}

// RunWatch starts the watch loop and blocks until SIGINT/SIGTERM.
func RunWatch(args []string, out io.Writer) error {
	wa, err := parseWatchArgs(args)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Watching %s vs %s (interval: %s)...\n", wa.LeftPath, wa.RightPath, wa.Interval)

	stop := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stop)
	}()

	opts := diff.WatchOptions{
		Interval: wa.Interval,
		OnChange: func(result diff.Result) {
			fmt.Fprintln(out, "--- change detected ---")
			diff.PrintFormatted(out, result, wa.Format)
		},
		OnError: func(e error) {
			fmt.Fprintf(out, "error: %v\n", e)
		},
	}

	diff.Watch(wa.LeftPath, wa.RightPath, opts, stop)
	fmt.Fprintln(out, "watch stopped")
	return nil
}

// matchFlag checks if arg matches --name=value and sets val.
func matchFlag(arg, name string, val *string) bool {
	prefix := name + "="
	if len(arg) > len(prefix) && arg[:len(prefix)] == prefix {
		*val = arg[len(prefix):]
		return true
	}
	return false
}
