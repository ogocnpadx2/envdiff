package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempEnvWatch(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	left := writeTempEnvWatch(t, dir, "left.env", "KEY=value1\n")
	right := writeTempEnvWatch(t, dir, "right.env", "KEY=value1\n")

	changed := make(chan Result, 1)
	stop := make(chan struct{})

	opts := WatchOptions{
		Interval: 50 * time.Millisecond,
		OnChange: func(r Result) {
			changed <- r
		},
	}

	go Watch(left, right, opts, stop)

	// Modify the right file after a short delay
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(right, []byte("KEY=value2\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case result := <-changed:
		if len(result.Mismatched) == 0 {
			t.Errorf("expected mismatch after file change, got none")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for change notification")
	}
	close(stop)
}

func TestWatch_NoSpuriousFire(t *testing.T) {
	dir := t.TempDir()
	left := writeTempEnvWatch(t, dir, "left.env", "KEY=same\n")
	right := writeTempEnvWatch(t, dir, "right.env", "KEY=same\n")

	fired := make(chan struct{}, 1)
	stop := make(chan struct{})

	opts := WatchOptions{
		Interval: 40 * time.Millisecond,
		OnChange: func(_ Result) {
			fired <- struct{}{}
		},
	}

	go Watch(left, right, opts, stop)
	time.Sleep(200 * time.Millisecond)
	close(stop)

	if len(fired) > 0 {
		t.Error("OnChange fired when files did not change")
	}
}

func TestFileHash_Consistent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.env")
	if err := os.WriteFile(p, []byte("A=1\nB=2\n"), 0644); err != nil {
		t.Fatal(err)
	}
	h1, err := fileHash(p)
	if err != nil {
		t.Fatalf("fileHash: %v", err)
	}
	h2, _ := fileHash(p)
	if h1 != h2 {
		t.Errorf("hash not consistent: %q vs %q", h1, h2)
	}
}
