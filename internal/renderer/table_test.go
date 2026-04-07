package renderer

import (
	"regexp"
	"strings"
	"testing"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleWidth(s string) int {
	plain := ansiPattern.ReplaceAllString(s, "")
	return len([]rune(plain))
}

func TestRenderSimpleTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "Name") {
		t.Fatalf("expected 'Name' in table output, got: %q", out)
	}
	if !strings.Contains(out, "Alice") {
		t.Fatalf("expected 'Alice' in table output, got: %q", out)
	}
	if !strings.Contains(out, "Bob") {
		t.Fatalf("expected 'Bob' in table output, got: %q", out)
	}
	if !strings.Contains(out, "─") {
		t.Fatalf("expected box-drawing borders in table output, got: %q", out)
	}
}

func TestRenderTableWithAlignment(t *testing.T) {
	input := "| Left | Center | Right |\n|:-----|:------:|------:|\n| a | b | c |"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "Left") {
		t.Fatalf("expected 'Left' in table output, got: %q", out)
	}
}

func TestRenderTableUnicodeCellsKeepBordersAligned(t *testing.T) {
	input := "| App | MCP Tier | Notes |\n|---|---|---|\n| OpenLM Chat UI | N/A — is the UI | First-party, core OS component |\n| GNOME Settings | Tier 2 | System settings. No MCP needed — system actions handled by System Control Layer. |"
	out := renderMarkdown(t, input)

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	var widths []int
	for _, line := range lines {
		if strings.Contains(line, "│") || strings.Contains(line, "┌") || strings.Contains(line, "└") || strings.Contains(line, "├") {
			widths = append(widths, visibleWidth(line))
		}
	}
	if len(widths) < 4 {
		t.Fatalf("expected table lines in output, got: %q", out)
	}
	want := widths[0]
	for i, got := range widths[1:] {
		if got != want {
			t.Fatalf("expected aligned table widths, line %d width=%d want=%d output=%q", i+2, got, want, out)
		}
	}
}
