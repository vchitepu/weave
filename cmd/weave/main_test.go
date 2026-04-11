package main

import (
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
)

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

	got := fileSeparator("notes.md", 80, th)

	if !strings.HasSuffix(got, "\n\n") {
		t.Fatalf("fileSeparator() should end with double newline, got: %q", got)
	}
}
