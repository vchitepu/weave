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

func TestRenderTableRespectsRendererWidth(t *testing.T) {
	input := "| Use case | Tool | Notes |\n|---|---|---|\n| OS and OpenLM system components | apt/deb | Managed by OpenLM, updated via standard Ubuntu apt repos + OpenLM apt repo |\n| User-facing apps | Flatpak | Via Flathub + OpenLM overlay repo |\n| Local LLM models | Ollama | Managed by Ollama daemon, not apt or Flatpak |\n| Python dependencies (LiteLLM) | pip (managed internally) | Isolated via virtualenv, not exposed to user |"
	out := renderMarkdown(t, input)

	max := 0
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if w := visibleWidth(line); w > max {
			max = w
		}
	}

	if max > 80 {
		t.Fatalf("expected table output max width <= 80, got %d output=%q", max, out)
	}
}
