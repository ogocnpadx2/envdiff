package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Annotation holds a note attached to a specific key.
type Annotation struct {
	Key     string
	Message string
}

// AnnotatedResult wraps a Result with per-key annotations.
type AnnotatedResult struct {
	Result      Result
	Annotations []Annotation
}

// AnnotateOptions controls which conditions produce annotations.
type AnnotateOptions struct {
	NoteEmpty    bool // annotate keys whose value is empty string
	NoteURL      bool // annotate keys whose value looks like a URL
	NotePlaceholder bool // annotate keys whose value looks like a placeholder
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		NoteEmpty:       true,
		NoteURL:         true,
		NotePlaceholder: true,
	}
}

// Annotate inspects the merged key/value pairs from both env files and
// attaches human-readable notes to interesting keys.
func Annotate(r Result, leftEnv, rightEnv map[string]string, opts AnnotateOptions) AnnotatedResult {
	var annotations []Annotation

	allKeys := unionKeys(leftEnv, rightEnv)
	sort.Strings(allKeys)

	for _, key := range allKeys {
		lv, lok := leftEnv[key]
		rv, rok := rightEnv[key]

		if opts.NoteEmpty {
			if lok && lv == "" {
				annotations = append(annotations, Annotation{Key: key, Message: "left value is empty"})
			}
			if rok && rv == "" {
				annotations = append(annotations, Annotation{Key: key, Message: "right value is empty"})
			}
		}

		if opts.NoteURL {
			if lok && isURL(lv) {
				annotations = append(annotations, Annotation{Key: key, Message: fmt.Sprintf("left value is a URL (%s)", lv)})
			}
			if rok && isURL(rv) && rv != lv {
				annotations = append(annotations, Annotation{Key: key, Message: fmt.Sprintf("right value is a URL (%s)", rv)})
			}
		}

		if opts.NotePlaceholder {
			if lok && isPlaceholder(lv) {
				annotations = append(annotations, Annotation{Key: key, Message: "left value looks like a placeholder"})
			}
			if rok && isPlaceholder(rv) {
				annotations = append(annotations, Annotation{Key: key, Message: "right value looks like a placeholder"})
			}
		}
	}

	return AnnotatedResult{Result: r, Annotations: annotations}
}

func isURL(v string) bool {
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}

func isPlaceholder(v string) bool {
	lower := strings.ToLower(v)
	for _, p := range []string{"todo", "changeme", "replace_me", "<your", "your_", "xxx"} {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func unionKeys(a, b map[string]string) []string {
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
