package diff

import (
	"testing"
)

func TestParseAliasMap_Valid(t *testing.T) {
	entries := []string{
		"DATABASE_URL=DB_URL,POSTGRES_URL",
		"SECRET_KEY=APP_SECRET",
	}
	am, err := ParseAliasMap(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(am["DATABASE_URL"]) != 2 {
		t.Errorf("expected 2 aliases for DATABASE_URL, got %d", len(am["DATABASE_URL"]))
	}
	if am["SECRET_KEY"][0] != "APP_SECRET" {
		t.Errorf("expected APP_SECRET alias, got %q", am["SECRET_KEY"][0])
	}
}

func TestParseAliasMap_Invalid(t *testing.T) {
	cases := []string{
		"NO_EQUALS",
		"=MISSING_CANONICAL",
		"MISSING_ALIAS=",
	}
	for _, c := range cases {
		_, err := ParseAliasMap([]string{c})
		if err == nil {
			t.Errorf("expected error for entry %q, got nil", c)
		}
	}
}

func TestResolveAliases_DirectMatch(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://localhost/db"}
	am := AliasMap{"DATABASE_URL": {"DB_URL", "POSTGRES_URL"}}
	res := ResolveAliases(env, am)
	if res.Resolved["DATABASE_URL"] != "postgres://localhost/db" {
		t.Errorf("expected direct value, got %q", res.Resolved["DATABASE_URL"])
	}
	if _, ok := res.UsedAlias["DATABASE_URL"]; ok {
		t.Error("should not record alias when canonical key present")
	}
}

func TestResolveAliases_FallbackToAlias(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost/db"}
	am := AliasMap{"DATABASE_URL": {"DB_URL", "POSTGRES_URL"}}
	res := ResolveAliases(env, am)
	if res.Resolved["DATABASE_URL"] != "postgres://localhost/db" {
		t.Errorf("expected alias-resolved value, got %q", res.Resolved["DATABASE_URL"])
	}
	if res.UsedAlias["DATABASE_URL"] != "DB_URL" {
		t.Errorf("expected UsedAlias DB_URL, got %q", res.UsedAlias["DATABASE_URL"])
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestResolveAliases_Unresolved(t *testing.T) {
	env := map[string]string{"UNRELATED": "value"}
	am := AliasMap{"DATABASE_URL": {"DB_URL"}}
	res := ResolveAliases(env, am)
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL unresolved, got %v", res.Unresolved)
	}
}

func TestResolveAliases_EmptyAliasMap(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	res := ResolveAliases(env, AliasMap{})
	if len(res.Resolved) != 0 || len(res.Unresolved) != 0 {
		t.Error("expected empty result for empty alias map")
	}
}
