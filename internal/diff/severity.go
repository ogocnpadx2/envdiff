package diff

// Severity represents the importance level of a diff finding.
type Severity int

const (
	SeverityInfo    Severity = iota // key exists but values match
	SeverityWarning                 // key missing in one file
	SeverityCritical                // key present in both but values differ
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// SeverityEntry pairs a key with its computed severity.
type SeverityEntry struct {
	Key      string
	Severity Severity
	Reason   string
}

// ClassifyResult inspects a Result and returns a slice of SeverityEntry
// describing the severity of each finding.
func ClassifyResult(r Result) []SeverityEntry {
	var entries []SeverityEntry

	for _, key := range r.MissingInRight {
		entries = append(entries, SeverityEntry{
			Key:      key,
			Severity: SeverityWarning,
			Reason:   "missing in right file",
		})
	}

	for _, key := range r.MissingInLeft {
		entries = append(entries, SeverityEntry{
			Key:      key,
			Severity: SeverityWarning,
			Reason:   "missing in left file",
		})
	}

	for _, m := range r.Mismatched {
		entries = append(entries, SeverityEntry{
			Key:      m.Key,
			Severity: SeverityCritical,
			Reason:   "value mismatch",
		})
	}

	return entries
}

// MaxSeverity returns the highest severity found in a Result.
// Returns SeverityInfo if there are no findings.
func MaxSeverity(r Result) Severity {
	entries := ClassifyResult(r)
	max := SeverityInfo
	for _, e := range entries {
		if e.Severity > max {
			max = e.Severity
		}
	}
	return max
}
