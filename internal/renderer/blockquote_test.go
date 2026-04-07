package renderer

import (
	"strings"
	"testing"
)

func TestRenderBlockquote(t *testing.T) {
	out := renderMarkdown(t, "> This is a quote")
	if !strings.Contains(out, "This is a quote") {
		t.Fatalf("expected quote text in output, got: %q", out)
	}
	if !strings.Contains(out, "│") {
		t.Fatalf("expected blockquote bar '│' in output, got: %q", out)
	}
}

func TestRenderNestedBlockquote(t *testing.T) {
	out := renderMarkdown(t, "> Outer\n>> Inner")
	if !strings.Contains(out, "Outer") {
		t.Fatalf("expected outer quote text, got: %q", out)
	}
	if !strings.Contains(out, "Inner") {
		t.Fatalf("expected inner quote text, got: %q", out)
	}
}

func TestRenderMultiLineBlockquote(t *testing.T) {
	out := renderMarkdown(t, "> Line one\n> Line two\n> Line three")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	barCount := 0
	for _, line := range lines {
		if strings.Contains(line, "│") {
			barCount++
		}
	}
	if barCount < 3 {
		t.Fatalf("expected blockquote bar on at least 3 lines, got %d bars in output: %q", barCount, out)
	}
}

func TestRenderLongBlockquoteWrapsWithBarOnEachLine(t *testing.T) {
	out := renderMarkdown(t, "> This is a very long blockquote line that should wrap across multiple rendered lines and keep the quote bar visible on each wrapped line")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	barCount := 0
	for _, line := range lines {
		if strings.Contains(line, "│") {
			barCount++
		}
	}

	if barCount < 2 {
		t.Fatalf("expected wrapped blockquote with bar on multiple lines, got: %q", out)
	}
}
