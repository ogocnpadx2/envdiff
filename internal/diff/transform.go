package diff

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a key-value pair.
type TransformFunc func(key, value string) (string, string, error)

// TransformOptions controls which transformations are applied.
type TransformOptions struct {
	TrimValues    bool
	LowercaseKeys bool
	UppercaseKeys bool
	PrefixKeys    string
	StripPrefix   string
}

// ParseTransformOptions builds TransformOptions from CLI-style flag strings.
func ParseTransformOptions(flags []string) (TransformOptions, error) {
	opts := TransformOptions{}
	for _, f := range flags {
		switch {
		case f == "trim":
			opts.TrimValues = true
		case f == "lowercase-keys":
			opts.LowercaseKeys = true
		case f == "uppercase-keys":
			opts.UppercaseKeys = true
		case strings.HasPrefix(f, "prefix="):
			opts.PrefixKeys = strings.TrimPrefix(f, "prefix=")
		case strings.HasPrefix(f, "strip-prefix="):
			opts.StripPrefix = strings.TrimPrefix(f, "strip-prefix=")
		default:
			return opts, fmt.Errorf("unknown transform option: %q", f)
		}
	}
	if opts.LowercaseKeys && opts.UppercaseKeys {
		return opts, fmt.Errorf("cannot combine lowercase-keys and uppercase-keys")
	}
	return opts, nil
}

// TransformEnv applies the given TransformOptions to a parsed env map,
// returning a new map with transformed keys and values.
func TransformEnv(env map[string]string, opts TransformOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.StripPrefix != "" {
			k = strings.TrimPrefix(k, opts.StripPrefix)
		}
		if opts.PrefixKeys != "" {
			k = opts.PrefixKeys + k
		}
		if opts.LowercaseKeys {
			k = strings.ToLower(k)
		} else if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		out[k] = v
	}
	return out
}
