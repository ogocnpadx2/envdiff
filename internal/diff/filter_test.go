package diff

import (
	"testing"
)

func baseResult() Result {
	return Result{
		MissingInRight: []string{"DB_HOST", "APP_SECRET"},
		MissingInLeft:  []string{"REDIS_URL"},
		Mismatched: []Mismatch{
			{Key: "APP_ENV", LeftVal: "development", RightVal: "production"},
			{Key: "DB_PORT", LeftVal: "5432", RightVal: "3306"},
		},
	}
}

func TestFilterResult_NoOptions(t *testing.T) {
	r := FilterResult(baseResult(), FilterOptions{})
	if len(r.MissingInRight) != 2 || len(r.MissingInLeft) != 1 || len(r.Mismatched) != 2 {
		t.Errorf("expected full result, got %+v", r)
	}
}

func TestFilterResult_OnlyMissing(t *testing.T) {
	r := FilterResult(baseResult(), FilterOptions{OnlyMissing: true})
	if len(r.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(r.MissingInRight))
	}
	if len(r.MissingInLeft) != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", len(r.MissingInLeft))
	}
	if len(r.Mismatched) != 0 {
		t.Errorf("expected 0 Mismatched, got %d", len(r.Mismatched))
	}
}

func TestFilterResult_OnlyMismatched(t *testing.T) {
	r := FilterResult(baseResult(), FilterOptions{OnlyMismatched: true})
	if len(r.Mismatched) != 2 {
		t.Errorf("expected 2 Mismatched, got %d", len(r.Mismatched))
	}
	if len(r.MissingInRight) != 0 || len(r.MissingInLeft) != 0 {
		t.Errorf("expected no missing keys, got right=%d left=%d", len(r.MissingInRight), len(r.MissingInLeft))
	}
}

func TestFilterResult_KeyPrefix(t *testing.T) {
	r := FilterResult(baseResult(), FilterOptions{KeyPrefix: "DB_"})
	if len(r.MissingInRight) != 1 || r.MissingInRight[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in MissingInRight, got %v", r.MissingInRight)
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT in Mismatched, got %v", r.Mismatched)
	}
	if len(r.MissingInLeft) != 0 {
		t.Errorf("expected no MissingInLeft, got %v", r.MissingInLeft)
	}
}

func TestFilterResult_PrefixCaseInsensitive(t *testing.T) {
	r := FilterResult(baseResult(), FilterOptions{KeyPrefix: "app_"})
	if len(r.MissingInRight) != 1 || r.MissingInRight[0] != "APP_SECRET" {
		t.Errorf("expected APP_SECRET, got %v", r.MissingInRight)
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV mismatch, got %v", r.Mismatched)
	}
}
