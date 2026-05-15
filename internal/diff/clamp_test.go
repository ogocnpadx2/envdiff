package diff

import (
	"testing"
)

func fptr(f float64) *float64 { return &f }

func TestClampEnv_NoOp_WhenNoBounds(t *testing.T) {
	env := map[string]string{"PORT": "8080", "TIMEOUT": "30"}
	report, err := ClampEnv(env, DefaultClampOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(report.Violations))
	}
	if report.Output["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", report.Output["PORT"])
	}
}

func TestClampEnv_ClampsMin(t *testing.T) {
	env := map[string]string{"WORKERS": "0", "TIMEOUT": "5"}
	opts := ClampOptions{Min: fptr(1)}
	report, err := ClampEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Output["WORKERS"] != "1" {
		t.Errorf("expected WORKERS=1, got %s", report.Output["WORKERS"])
	}
	if report.Output["TIMEOUT"] != "5" {
		t.Errorf("expected TIMEOUT=5, got %s", report.Output["TIMEOUT"])
	}
	if len(report.Violations) != 1 || report.Violations[0].Key != "WORKERS" {
		t.Errorf("expected 1 violation for WORKERS, got %+v", report.Violations)
	}
}

func TestClampEnv_ClampsMax(t *testing.T) {
	env := map[string]string{"PORT": "99999"}
	opts := ClampOptions{Max: fptr(65535)}
	report, err := ClampEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Output["PORT"] != "65535" {
		t.Errorf("expected PORT=65535, got %s", report.Output["PORT"])
	}
	if len(report.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(report.Violations))
	}
}

func TestClampEnv_SkipsNonNumeric(t *testing.T) {
	env := map[string]string{"NAME": "alice", "PORT": "3000"}
	opts := ClampOptions{Min: fptr(1), Max: fptr(9000)}
	report, err := ClampEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Output["NAME"] != "alice" {
		t.Errorf("expected NAME=alice, got %s", report.Output["NAME"])
	}
	if len(report.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(report.Violations))
	}
}

func TestClampEnv_StrictRejectsNonNumeric(t *testing.T) {
	env := map[string]string{"PORT": "notanumber"}
	opts := ClampOptions{Min: fptr(0), Keys: []string{"PORT"}, Strict: true}
	_, err := ClampEnv(env, opts)
	if err == nil {
		t.Fatal("expected error for non-numeric value in strict mode")
	}
}

func TestClampEnv_KeyFilter_OnlyTargeted(t *testing.T) {
	env := map[string]string{"PORT": "0", "RETRIES": "0"}
	opts := ClampOptions{Min: fptr(1), Keys: []string{"PORT"}}
	report, err := ClampEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Output["PORT"] != "1" {
		t.Errorf("expected PORT=1, got %s", report.Output["PORT"])
	}
	if report.Output["RETRIES"] != "0" {
		t.Errorf("expected RETRIES=0 (untargeted), got %s", report.Output["RETRIES"])
	}
	if len(report.Violations) != 1 || report.Violations[0].Key != "PORT" {
		t.Errorf("expected 1 violation for PORT only, got %+v", report.Violations)
	}
}
