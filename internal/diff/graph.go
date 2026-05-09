package diff

import "sort"

// DependencyGraph maps each key to the set of keys it references via interpolation.
type DependencyGraph map[string][]string

// GraphEntry describes a single node in the dependency graph.
type GraphEntry struct {
	Key  string
	Deps []string
}

// BuildDependencyGraph analyses an env map and returns a graph of which keys
// reference other keys (via $VAR or ${VAR} syntax).
func BuildDependencyGraph(env map[string]string) DependencyGraph {
	graph := make(DependencyGraph, len(env))
	for key, val := range env {
		deps := extractRefs(val)
		// Only include deps that are themselves defined keys.
		var filtered []string
		for _, d := range deps {
			if _, ok := env[d]; ok {
				filtered = append(filtered, d)
			}
		}
		sort.Strings(filtered)
		graph[key] = filtered
	}
	return graph
}

// SortedEntries returns graph entries in deterministic key order.
func (g DependencyGraph) SortedEntries() []GraphEntry {
	keys := make([]string, 0, len(g))
	for k := range g {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]GraphEntry, 0, len(keys))
	for _, k := range keys {
		out = append(out, GraphEntry{Key: k, Deps: g[k]})
	}
	return out
}

// CyclicKeys returns any keys that participate in a dependency cycle.
func (g DependencyGraph) CyclicKeys() []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles []string

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		for _, dep := range g[node] {
			if !visited[dep] {
				if dfs(dep) {
					return true
				}
			} else if recStack[dep] {
				cycles = append(cycles, node)
				return true
			}
		}
		recStack[node] = false
		return false
	}

	for node := range g {
		if !visited[node] {
			dfs(node)
		}
	}
	sort.Strings(cycles)
	return cycles
}

// extractRefs returns all variable names referenced in a value string.
func extractRefs(val string) []string {
	var refs []string
	seen := make(map[string]bool)
	i := 0
	for i < len(val) {
		if val[i] == '$' && i+1 < len(val) {
			i++
			var name string
			if val[i] == '{' {
				i++
				start := i
				for i < len(val) && val[i] != '}' {
					i++
				}
				name = val[start:i]
				if i < len(val) {
					i++
				}
			} else {
				start := i
				for i < len(val) && isIdentChar(val[i]) {
					i++
				}
				name = val[start:i]
			}
			if name != "" && !seen[name] {
				seen[name] = true
				refs = append(refs, name)
			}
		} else {
			i++
		}
	}
	return refs
}

func isIdentChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}
