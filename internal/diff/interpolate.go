package diff

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation is performed.
type InterpolateOptions struct {
	// FallbackToOS allows falling back to OS environment variables when a key
	// is not found in the provided env map.
	FallbackToOS bool
	// FailOnMissing returns an error if a referenced variable cannot be resolved.
	FailOnMissing bool
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		FallbackToOS:  false,
		FailOnMissing: false,
	}
}

// InterpolateResult holds the resolved env map and any unresolved variable names.
type InterpolateResult struct {
	Resolved   map[string]string
	Unresolved []string
}

// InterpolateEnv expands variable references within env values using other
// values in the same map (and optionally the OS environment).
func InterpolateEnv(env map[string]string, opts InterpolateOptions) (InterpolateResult, error) {
	resolved := make(map[string]string, len(env))
	unresolvedSet := map[string]struct{}{}

	for key, value := range env {
		expanded, missing, err := expandValue(value, env, opts)
		if err != nil {
			return InterpolateResult{}, fmt.Errorf("key %q: %w", key, err)
		}
		resolved[key] = expanded
		for _, m := range missing {
			unresolvedSet[m] = struct{}{}
		}
	}

	unresolved := make([]string, 0, len(unresolvedSet))
	for k := range unresolvedSet {
		unresolved = append(unresolved, k)
	}
	sortStrings(unresolved)

	return InterpolateResult{Resolved: resolved, Unresolved: unresolved}, nil
}

func expandValue(value string, env map[string]string, opts InterpolateOptions) (string, []string, error) {
	var missing []string
	var expandErr error

	result := interpolatePattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		submatches := interpolatePattern.FindStringSubmatch(match)
		varName := submatches[1]
		if varName == "" {
			varName = submatches[2]
		}

		if v, ok := env[varName]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(varName); ok {
				return v
			}
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("unresolved variable: %s", varName)
			return match
		}
		missing = append(missing, varName)
		return strings.TrimSpace(match) // leave as-is
	})

	if expandErr != nil {
		return "", nil, expandErr
	}
	return result, missing, nil
}
