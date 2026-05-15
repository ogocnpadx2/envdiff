package diff

import (
	"fmt"
	"sort"
	"strconv"
)

// ClampOptions controls how numeric env values are clamped.
type ClampOptions struct {
	Min    *float64
	Max    *float64
	Keys   []string // if empty, apply to all numeric values
	Strict bool     // if true, return error on non-numeric values in targeted keys
}

// ClampViolation records a key whose value was clamped or rejected.
type ClampViolation struct {
	Key      string
	Original string
	Result   string
	Reason   string
}

// ClampReport is the result of a clamp operation.
type ClampReport struct {
	Output     map[string]string
	Violations []ClampViolation
}

// DefaultClampOptions returns a ClampOptions with no bounds set.
func DefaultClampOptions() ClampOptions {
	return ClampOptions{}
}

// ClampEnv applies numeric clamping to env values according to opts.
func ClampEnv(env map[string]string, opts ClampOptions) (ClampReport, error) {
	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[k] = true
	}

	out := make(map[string]string, len(env))
	var violations []ClampViolation

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		targeted := len(targetSet) == 0 || targetSet[k]

		if !targeted {
			out[k] = v
			continue
		}

		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			if opts.Strict && targeted && len(targetSet) > 0 {
				return ClampReport{}, fmt.Errorf("key %q: value %q is not numeric", k, v)
			}
			out[k] = v
			continue
		}

		clamped := f
		reason := ""
		if opts.Min != nil && f < *opts.Min {
			clamped = *opts.Min
			reason = fmt.Sprintf("below min (%.6g)", *opts.Min)
		} else if opts.Max != nil && f > *opts.Max {
			clamped = *opts.Max
			reason = fmt.Sprintf("above max (%.6g)", *opts.Max)
		}

		resultStr := strconv.FormatFloat(clamped, 'f', -1, 64)
		out[k] = resultStr

		if reason != "" {
			violations = append(violations, ClampViolation{
				Key:      k,
				Original: v,
				Result:   resultStr,
				Reason:   reason,
			})
		}
	}

	return ClampReport{Output: out, Violations: violations}, nil
}
