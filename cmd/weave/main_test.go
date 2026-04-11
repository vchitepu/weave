package main

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleWidth(s string) int {
	plain := ansiPattern.ReplaceAllString(s, "")
	return len([]rune(plain))
}

func TestNormalizeWidth_AutoWidthCappedAt120(t *testing.T) {
	if got := normalizeWidth(221, true); got != 120 {
		t.Fatalf("normalizeWidth(221, auto=true) = %d, want 120", got)
	}
}

func TestNormalizeWidth_ExplicitWidthNotCapped(t *testing.T) {
	if got := normalizeWidth(400, false); got != 400 {
		t.Fatalf("normalizeWidth(400, auto=false) = %d, want 400", got)
	}
}

func TestNormalizeWidth_Minimum20(t *testing.T) {
	if got := normalizeWidth(10, false); got != 20 {
		t.Fatalf("normalizeWidth(10, auto=false) = %d, want 20", got)
	}
}

func TestFileSeparator_ContainsFilename(t *testing.T) {
	th := theme.DarkTheme()

	got := fileSeparator("notes.md", 80, th)

	if !strings.Contains(got, "notes.md") {
		t.Fatalf("fileSeparator() output missing filename: %q", got)
	}
}

func TestFileSeparator_ContainsRule(t *testing.T) {
	th := theme.DarkTheme()

	got := fileSeparator("notes.md", 80, th)

	if !strings.Contains(got, "─") {
		t.Fatalf("fileSeparator() output missing horizontal rule: %q", got)
	}
}

func TestFileSeparator_EndsWithNewline(t *testing.T) {
	th := theme.DarkTheme()

	sep := fileSeparator("notes.md", 80, th)

	if !strings.HasSuffix(sep, "\n") {
		t.Fatalf("fileSeparator() should end with newline, got: %q", sep)
	}
}

func TestFileSeparator_RuleDoesNotExceedTerminalWidth(t *testing.T) {
	th := theme.DarkTheme()

	sep := fileSeparator("notes.md", 20, th)
	lines := strings.Split(strings.TrimRight(sep, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("fileSeparator() unexpected output: %q", sep)
	}

	ruleLine := lines[1]
	if got := visibleWidth(ruleLine); got > 20 {
		t.Fatalf("fileSeparator() rule line width = %d, want <= 20; output=%q", got, sep)
	}
}

func TestRenderFile_ValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "weave-*.md")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString("# Hello\n"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}

	md := buildMarkdown(theme.DarkTheme(), 80)

	got, err := renderFile(tmpFile.Name(), md)
	if err != nil {
		t.Fatalf("renderFile() error = %v", err)
	}
	if strings.TrimSpace(got) == "" {
		t.Fatalf("renderFile() output should be non-empty")
	}
	if !strings.Contains(got, "Hello") {
		t.Fatalf("renderFile() output missing expected content: %q", got)
	}
}

func TestRenderFile_MissingFile(t *testing.T) {
	missingPath := "/tmp/does-not-exist-weave-test.md"
	md := buildMarkdown(theme.DarkTheme(), 80)

	_, err := renderFile(missingPath, md)
	if err == nil {
		t.Fatalf("renderFile() error = nil, want missing-file error")
	}

	want := "weave: no such file: " + missingPath
	if got := err.Error(); got != want {
		t.Fatalf("renderFile() error = %q, want %q", got, want)
	}
}

func TestMultiFileOutput(t *testing.T) {
	tmp1, err := os.CreateTemp(t.TempDir(), "weave-multi1-*.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmp1.WriteString("# First\n\nAlpha content\n"); err != nil {
		t.Fatal(err)
	}
	if err := tmp1.Close(); err != nil {
		t.Fatal(err)
	}

	tmp2, err := os.CreateTemp(t.TempDir(), "weave-multi2-*.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmp2.WriteString("# Second\n\nBeta content\n"); err != nil {
		t.Fatal(err)
	}
	if err := tmp2.Close(); err != nil {
		t.Fatal(err)
	}

	th := theme.DarkTheme()
	md := buildMarkdown(th, 80)

	var combined strings.Builder
	paths := []string{tmp1.Name(), tmp2.Name()}
	for i, path := range paths {
		if i > 0 {
			combined.WriteString(fileSeparator(path, 80, th))
		}

		rendered, err := renderFile(path, md)
		if err != nil {
			t.Fatalf("renderFile(%s): %v", path, err)
		}
		combined.WriteString(rendered)
	}

	output := combined.String()

	if !strings.Contains(output, "First") {
		t.Error("output should contain 'First'")
	}
	if !strings.Contains(output, "Alpha") {
		t.Error("output should contain 'Alpha'")
	}
	if !strings.Contains(output, "Second") {
		t.Error("output should contain 'Second'")
	}
	if !strings.Contains(output, "Beta") {
		t.Error("output should contain 'Beta'")
	}
	if !strings.Contains(output, "─") {
		t.Error("output should contain a rule character in the separator")
	}

	secondFilename := tmp2.Name()
	if !strings.Contains(output, secondFilename) {
		t.Errorf("output should contain second filename %q in separator", secondFilename)
	}
}
