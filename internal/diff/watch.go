package diff

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the file watcher behaviour.
type WatchOptions struct {
	Interval  time.Duration
	OnChange  func(result Result)
	OnError   func(err error)
}

// fileHash returns an MD5 hash string for the given file path.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls two .env files at the given interval and calls opts.OnChange
// whenever either file changes. It blocks until stop is closed.
func Watch(leftPath, rightPath string, opts WatchOptions, stop <-chan struct{}) {
	if opts.Interval <= 0 {
		opts.Interval = 2 * time.Second
	}

	prevLeft, _ := fileHash(leftPath)
	prevRight, _ := fileHash(rightPath)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			curLeft, errL := fileHash(leftPath)
			curRight, errR := fileHash(rightPath)
			if errL != nil || errR != nil {
				if opts.OnError != nil {
					if errL != nil {
						opts.OnError(errL)
					} else {
						opts.OnError(errR)
					}
				}
				continue
			}
			if curLeft != prevLeft || curRight != prevRight {
				prevLeft = curLeft
				prevRight = curRight
				result, err := Compare(leftPath, rightPath)
				if err != nil {
					if opts.OnError != nil {
						opts.OnError(err)
					}
					continue
				}
				if opts.OnChange != nil {
					opts.OnChange(result)
				}
			}
		}
	}
}
