package diff_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatch_ErrorOnMissingFile(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.env")
	right := filepath.Join(dir, "right.env")

	// Only create left; right is missing
	if err := os.WriteFile(left, []byte("K=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	stop := make(chan struct{})

	opts := WatchOptions{
		Interval: 40 * time.Millisecond,
		OnError: func(e error) {
			select {
			case errCh <- e:
			default:
			}
		},
	}

	go Watch(left, right, opts, stop)

	select {
	case err := <-errCh:
		if err == nil {
			t.Error("expected non-nil error")
		}
	case <-time.After(400 * time.Millisecond):
		t.Error("timed out waiting for error callback")
	}
	close(stop)
}

func TestWatch_StopsCleanly(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "l.env")
	right := filepath.Join(dir, "r.env")
	if err := os.WriteFile(left, []byte("A=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("A=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		Watch(left, right, WatchOptions{Interval: 30 * time.Millisecond}, stop)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	close(stop)

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Error("Watch did not stop after stop channel closed")
	}
}
