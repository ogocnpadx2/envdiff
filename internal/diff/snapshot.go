package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the result of a diff at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	LeftFile  string    `json:"left_file"`
	RightFile string    `json:"right_file"`
	Result    Result    `json:"result"`
}

// SaveSnapshot writes a Snapshot to the given file path as JSON.
func SaveSnapshot(path, leftFile, rightFile string, result Result) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		LeftFile:  leftFile,
		RightFile: rightFile,
		Result:    result,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// LoadSnapshot reads and parses a Snapshot from the given file path.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: parse failed: %w", err)
	}
	return &snap, nil
}

// DiffSnapshots compares two snapshots and returns keys whose status changed.
func DiffSnapshots(before, after *Snapshot) SnapshotDelta {
	delta := SnapshotDelta{
		Before: before.Timestamp,
		After:  after.Timestamp,
	}

	beforeKeys := snapshotKeySet(before.Result)
	afterKeys := snapshotKeySet(after.Result)

	for k := range afterKeys {
		if _, existed := beforeKeys[k]; !existed {
			delta.NewIssues = append(delta.NewIssues, k)
		}
	}
	for k := range beforeKeys {
		if _, exists := afterKeys[k]; !exists {
			delta.ResolvedIssues = append(delta.ResolvedIssues, k)
		}
	}
	sortStrings(delta.NewIssues)
	sortStrings(delta.ResolvedIssues)
	return delta
}

// SnapshotDelta describes what changed between two snapshots.
type SnapshotDelta struct {
	Before         time.Time `json:"before"`
	After          time.Time `json:"after"`
	NewIssues      []string  `json:"new_issues"`
	ResolvedIssues []string  `json:"resolved_issues"`
}

func snapshotKeySet(r Result) map[string]struct{} {
	set := make(map[string]struct{})
	for _, k := range r.MissingInRight {
		set[k] = struct{}{}
	}
	for _, k := range r.MissingInLeft {
		set[k] = struct{}{}
	}
	for _, m := range r.Mismatched {
		set[m.Key] = struct{}{}
	}
	return set
}
