package diff

import (
	"testing"
)

func TestBuildAuditLog_AllUnchanged(t *testing.T) {
	left := map[string]string{"KEY": "val", "FOO": "bar"}
	right := map[string]string{"KEY": "val", "FOO": "bar"}
	log := BuildAuditLog("a.env", "b.env", left, right)

	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	for _, e := range log.Events {
		if e.Type != AuditUnchanged {
			t.Errorf("expected unchanged for %s, got %s", e.Key, e.Type)
		}
	}
}

func TestBuildAuditLog_Added(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"NEW_KEY": "value"}
	log := BuildAuditLog("a.env", "b.env", left, right)

	if len(log.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(log.Events))
	}
	e := log.Events[0]
	if e.Type != AuditAdded {
		t.Errorf("expected added, got %s", e.Type)
	}
	if e.NewValue != "value" {
		t.Errorf("expected new value 'value', got %q", e.NewValue)
	}
}

func TestBuildAuditLog_Removed(t *testing.T) {
	left := map[string]string{"OLD_KEY": "gone"}
	right := map[string]string{}
	log := BuildAuditLog("a.env", "b.env", left, right)

	if len(log.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(log.Events))
	}
	e := log.Events[0]
	if e.Type != AuditRemoved {
		t.Errorf("expected removed, got %s", e.Type)
	}
	if e.OldValue != "gone" {
		t.Errorf("expected old value 'gone', got %q", e.OldValue)
	}
}

func TestBuildAuditLog_Changed(t *testing.T) {
	left := map[string]string{"KEY": "old"}
	right := map[string]string{"KEY": "new"}
	log := BuildAuditLog("a.env", "b.env", left, right)

	if len(log.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(log.Events))
	}
	e := log.Events[0]
	if e.Type != AuditChanged {
		t.Errorf("expected changed, got %s", e.Type)
	}
	if e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected values: old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAuditLog_Summary(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old", "C": "same"}
	right := map[string]string{"B": "new", "C": "same", "D": "added"}
	log := BuildAuditLog("a.env", "b.env", left, right)

	summary := log.Summary()
	expected := "added=1 removed=1 changed=1 unchanged=1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestBuildAuditLog_SortedKeys(t *testing.T) {
	left := map[string]string{"ZZZ": "1", "AAA": "2"}
	right := map[string]string{"ZZZ": "1", "AAA": "2"}
	log := BuildAuditLog("a.env", "b.env", left, right)

	if log.Events[0].Key != "AAA" || log.Events[1].Key != "ZZZ" {
		t.Errorf("events not sorted: got %s, %s", log.Events[0].Key, log.Events[1].Key)
	}
}
