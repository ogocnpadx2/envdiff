package diff

import (
	"fmt"
	"sort"
	"time"
)

// AuditEventType describes what kind of audit event occurred.
type AuditEventType string

const (
	AuditAdded    AuditEventType = "added"
	AuditRemoved  AuditEventType = "removed"
	AuditChanged  AuditEventType = "changed"
	AuditUnchanged AuditEventType = "unchanged"
)

// AuditEvent represents a single key-level change between two env snapshots.
type AuditEvent struct {
	Key       string         `json:"key"`
	Type      AuditEventType `json:"type"`
	OldValue  string         `json:"old_value,omitempty"`
	NewValue  string         `json:"new_value,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// AuditLog is a collection of audit events with metadata.
type AuditLog struct {
	GeneratedAt time.Time    `json:"generated_at"`
	LeftFile    string       `json:"left_file"`
	RightFile   string       `json:"right_file"`
	Events      []AuditEvent `json:"events"`
}

// BuildAuditLog constructs an AuditLog from two parsed env maps.
func BuildAuditLog(leftFile, rightFile string, left, right map[string]string) AuditLog {
	now := time.Now().UTC()
	keys := unionAuditKeys(left, right)
	sort.Strings(keys)

	events := make([]AuditEvent, 0, len(keys))
	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]

		var ev AuditEvent
		ev.Key = k
		ev.Timestamp = now

		switch {
		case inLeft && !inRight:
			ev.Type = AuditRemoved
			ev.OldValue = lv
		case !inLeft && inRight:
			ev.Type = AuditAdded
			ev.NewValue = rv
		case lv != rv:
			ev.Type = AuditChanged
			ev.OldValue = lv
			ev.NewValue = rv
		default:
			ev.Type = AuditUnchanged
		}
		events = append(events, ev)
	}

	return AuditLog{
		GeneratedAt: now,
		LeftFile:    leftFile,
		RightFile:   rightFile,
		Events:      events,
	}
}

// Summary returns a human-readable summary line for the audit log.
func (a *AuditLog) Summary() string {
	var added, removed, changed, unchanged int
	for _, e := range a.Events {
		switch e.Type {
		case AuditAdded:
			added++
		case AuditRemoved:
			removed++
		case AuditChanged:
			changed++
		case AuditUnchanged:
			unchanged++
		}
	}
	return fmt.Sprintf("added=%d removed=%d changed=%d unchanged=%d", added, removed, changed, unchanged)
}

func unionAuditKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
