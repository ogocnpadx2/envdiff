package diff

import (
	"testing"
)

func TestAnnotate_NoAnnotations(t *testing.T) {
	left := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	right := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	r := Compare(left, right)
	ar := Annotate(r, left, right, DefaultAnnotateOptions())
	if len(ar.Annotations) != 0 {
		t.Errorf("expected no annotations, got %d", len(ar.Annotations))
	}
}

func TestAnnotate_EmptyValue(t *testing.T) {
	left := map[string]string{"SECRET": ""}
	right := map[string]string{"SECRET": "abc"}
	r := Compare(left, right)
	ar := Annotate(r, left, right, DefaultAnnotateOptions())
	found := false
	for _, a := range ar.Annotations {
		if a.Key == "SECRET" && a.Message == "left value is empty" {
			found = true
		}
	}
	if !found {
		t.Error("expected annotation for empty left value on SECRET")
	}
}

func TestAnnotate_URLValue(t *testing.T) {
	left := map[string]string{"API_URL": "https://api.example.com"}
	right := map[string]string{"API_URL": "https://api.example.com"}
	r := Compare(left, right)
	opts := DefaultAnnotateOptions()
	ar := Annotate(r, left, right, opts)
	found := false
	for _, a := range ar.Annotations {
		if a.Key == "API_URL" {
			found = true
		}
	}
	if !found {
		t.Error("expected URL annotation for API_URL")
	}
}

func TestAnnotate_Placeholder(t *testing.T) {
	left := map[string]string{"DB_PASS": "changeme"}
	right := map[string]string{"DB_PASS": "s3cr3t"}
	r := Compare(left, right)
	ar := Annotate(r, left, right, DefaultAnnotateOptions())
	found := false
	for _, a := range ar.Annotations {
		if a.Key == "DB_PASS" && a.Message == "left value looks like a placeholder" {
			found = true
		}
	}
	if !found {
		t.Error("expected placeholder annotation for DB_PASS left value")
	}
}

func TestAnnotate_DisabledOptions(t *testing.T) {
	left := map[string]string{"DB_PASS": "changeme", "ENDPOINT": "https://x.com", "KEY": ""}
	right := map[string]string{"DB_PASS": "changeme", "ENDPOINT": "https://x.com", "KEY": ""}
	r := Compare(left, right)
	opts := AnnotateOptions{NoteEmpty: false, NoteURL: false, NotePlaceholder: false}
	ar := Annotate(r, left, right, opts)
	if len(ar.Annotations) != 0 {
		t.Errorf("expected no annotations with all options disabled, got %d", len(ar.Annotations))
	}
}

func TestAnnotate_PreservesResult(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"B": "2"}
	r := Compare(left, right)
	ar := Annotate(r, left, right, DefaultAnnotateOptions())
	if len(ar.Result.MissingInRight) != 1 || ar.Result.MissingInRight[0] != "A" {
		t.Error("annotated result should preserve original diff result")
	}
}
