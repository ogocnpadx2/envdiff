package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// EncryptOptions controls how values are masked or hashed in output.
type EncryptOptions struct {
	// HashKeys is the list of substrings; matching keys will have values hashed.
	HashKeys []string
	// Prefix is prepended to hashed values in output.
	Prefix string
}

// DefaultEncryptOptions returns sensible defaults.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		HashKeys: []string{"secret", "password", "token", "key", "auth", "credential"},
		Prefix:   "sha256:",
	}
}

// hashValue returns a short SHA-256 hex digest of the input.
func hashValue(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])[:16]
}

// shouldHash returns true when the key matches any hash keyword.
func shouldHash(key string, keywords []string) bool {
	lower := strings.ToLower(key)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// EncryptEnv returns a copy of env with sensitive values replaced by hashes.
func EncryptEnv(env map[string]string, opts EncryptOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldHash(k, opts.HashKeys) {
			out[k] = fmt.Sprintf("%s%s", opts.Prefix, hashValue(v))
		} else {
			out[k] = v
		}
	}
	return out
}

// EncryptReport summarises which keys were hashed.
type EncryptReport struct {
	HashedKeys []string
}

// EncryptEnvWithReport hashes sensitive values and returns a report.
func EncryptEnvWithReport(env map[string]string, opts EncryptOptions) (map[string]string, EncryptReport) {
	out := make(map[string]string, len(env))
	var hashed []string
	for k, v := range env {
		if shouldHash(k, opts.HashKeys) {
			out[k] = fmt.Sprintf("%s%s", opts.Prefix, hashValue(v))
			hashed = append(hashed, k)
		} else {
			out[k] = v
		}
	}
	sort.Strings(hashed)
	return out, EncryptReport{HashedKeys: hashed}
}
