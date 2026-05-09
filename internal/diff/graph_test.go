package diff

import (
	"testing"
)

func TestBuildDependencyGraph_NoDeps(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	g := BuildDependencyGraph(env)
	if len(g["HOST"]) != 0 {
		t.Errorf("expected no deps for HOST, got %v", g["HOST"])
	}
}

func TestBuildDependencyGraph_SimpleDep(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	}
	g := BuildDependencyGraph(env)
	if len(g["DSN"]) != 1 || g["DSN"][0] != "HOST" {
		t.Errorf("expected DSN -> [HOST], got %v", g["DSN"])
	}
}

func TestBuildDependencyGraph_CurlyBrace(t *testing.T) {
	env := map[string]string{
		"USER": "admin",
		"URL":  "http://${USER}@example.com",
	}
	g := BuildDependencyGraph(env)
	if len(g["URL"]) != 1 || g["URL"][0] != "USER" {
		t.Errorf("expected URL -> [USER], got %v", g["URL"])
	}
}

func TestBuildDependencyGraph_UndefinedRefIgnored(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://$UNKNOWN_HOST/db",
	}
	g := BuildDependencyGraph(env)
	if len(g["DSN"]) != 0 {
		t.Errorf("expected undefined refs to be filtered, got %v", g["DSN"])
	}
}

func TestSortedEntries_Order(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	g := BuildDependencyGraph(env)
	entries := g.SortedEntries()
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestCyclicKeys_NoCycle(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "$A world",
	}
	g := BuildDependencyGraph(env)
	cycles := g.CyclicKeys()
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}

func TestCyclicKeys_DirectCycle(t *testing.T) {
	// Manually construct a cycle since the parser won't produce one naturally.
	g := DependencyGraph{
		"A": {"B"},
		"B": {"A"},
	}
	cycles := g.CyclicKeys()
	if len(cycles) == 0 {
		t.Error("expected cycle to be detected")
	}
}

func TestExtractRefs_MultipleSameRef(t *testing.T) {
	refs := extractRefs("$FOO/$FOO")
	if len(refs) != 1 || refs[0] != "FOO" {
		t.Errorf("expected deduplication, got %v", refs)
	}
}
