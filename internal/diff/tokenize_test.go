package diff

import (
	"testing"
)

func TestTokenizeValues_NoMismatch(t *testing.T) {
	left := map[string]string{"FEATURES": "a,b,c"}
	right := map[string]string{"FEATURES": "a,b,c"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestTokenizeValues_OnlyInLeft(t *testing.T) {
	left := map[string]string{"FEATURES": "a,b,c"}
	right := map[string]string{"FEATURES": "a,b"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0] != "c" {
		t.Errorf("expected OnlyInLeft=[c], got %v", r.OnlyInLeft)
	}
	if len(r.OnlyInRight) != 0 {
		t.Errorf("expected empty OnlyInRight, got %v", r.OnlyInRight)
	}
	if len(r.Shared) != 2 {
		t.Errorf("expected 2 shared tokens, got %v", r.Shared)
	}
}

func TestTokenizeValues_OnlyInRight(t *testing.T) {
	left := map[string]string{"SCOPES": "read"}
	right := map[string]string{"SCOPES": "read,write"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "write" {
		t.Errorf("expected OnlyInRight=[write], got %v", r.OnlyInRight)
	}
}

func TestTokenizeValues_CaseInsensitive(t *testing.T) {
	left := map[string]string{"FLAGS": "Alpha,Beta"}
	right := map[string]string{"FLAGS": "alpha,gamma"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.Shared) != 1 || r.Shared[0] != "alpha" {
		t.Errorf("expected shared=[alpha], got %v", r.Shared)
	}
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "gamma" {
		t.Errorf("expected OnlyInRight=[gamma], got %v", r.OnlyInRight)
	}
}

func TestTokenizeValues_MissingKeySkipped(t *testing.T) {
	left := map[string]string{"A": "x,y", "B": "1,2"}
	right := map[string]string{"A": "x,z"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result (only A), got %d", len(results))
	}
	if results[0].Key != "A" {
		t.Errorf("expected key A, got %s", results[0].Key)
	}
}

func TestTokenizeValues_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1,2", "A": "a,b"}
	right := map[string]string{"Z": "1,3", "A": "a,c"}
	results := TokenizeValues(left, right, DefaultTokenizeOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Key != "A" || results[1].Key != "Z" {
		t.Errorf("expected sorted keys [A Z], got [%s %s]", results[0].Key, results[1].Key)
	}
}

func TestTokenizeValues_CustomDelimiter(t *testing.T) {
	opts := TokenizeOptions{Delimiter: "|", Lowercase: false}
	left := map[string]string{"HOSTS": "host1|host2"}
	right := map[string]string{"HOSTS": "host1|host3"}
	results := TokenizeValues(left, right, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0] != "host2" {
		t.Errorf("expected OnlyInLeft=[host2], got %v", r.OnlyInLeft)
	}
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "host3" {
		t.Errorf("expected OnlyInRight=[host3], got %v", r.OnlyInRight)
	}
}
