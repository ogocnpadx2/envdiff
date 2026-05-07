package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempTransformEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestParseTransformArgs_MissingPath(t *testing.T) {
	_, err := parseTransformArgs([]string{})
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestParseTransformArgs_Defaults(t *testing.T) {
	a, err := parseTransformArgs([]string{"file.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.filePath != "file.env" {
		t.Errorf("expected file.env, got %q", a.filePath)
	}
	if a.output != "text" {
		t.Errorf("expected output=text, got %q", a.output)
	}
}

func TestParseTransformArgs_UnknownFlag(t *testing.T) {
	_, err := parseTransformArgs([]string{"file.env", "--unknown"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestRunTransform_TrimAndUppercase(t *testing.T) {
	path := writeTempTransformEnv(t, "my_key=  hello  \nother_key=world\n")
	var buf bytes.Buffer
	err := RunTransform([]string{path, "--transform=trim,uppercase-keys"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MY_KEY") {
		t.Errorf("expected MY_KEY in output, got:\n%s", out)
	}
	if strings.Contains(out, "  hello  ") {
		t.Errorf("expected trimmed value, got:\n%s", out)
	}
}

func TestRunTransform_DotEnvOutput(t *testing.T) {
	path := writeTempTransformEnv(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunTransform([]string{path, "--output=dotenv"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY=value") {
		t.Errorf("expected KEY=value in dotenv output, got:\n%s", out)
	}
}

func TestRunTransform_InvalidTransformOption(t *testing.T) {
	path := writeTempTransformEnv(t, "KEY=val\n")
	var buf bytes.Buffer
	err := RunTransform([]string{path, "--transform=bad-option"}, &buf)
	if err == nil {
		t.Error("expected error for invalid transform option")
	}
}
