package renderer

import (
	"strings"
	"testing"
)

func TestRenderH1(t *testing.T) {
	out := renderMarkdown(t, "# Hello")
	if !strings.Contains(out, "Hello") {
		t.Fatalf("expected H1 output to contain 'Hello', got: %q", out)
	}
	// H1 should have a horizontal rule beneath it (─)
	if !strings.Contains(out, "─") {
		t.Fatalf("expected H1 to have a rule beneath, got: %q", out)
	}
}

func TestRenderH2(t *testing.T) {
	out := renderMarkdown(t, "## Subtitle")
	if !strings.Contains(out, "Subtitle") {
		t.Fatalf("expected H2 output to contain 'Subtitle', got: %q", out)
	}
}

func TestRenderH3(t *testing.T) {
	out := renderMarkdown(t, "### Section")
	if !strings.Contains(out, "Section") {
		t.Fatalf("expected H3 output to contain 'Section', got: %q", out)
	}
}

func TestRenderH4UsesH3Style(t *testing.T) {
	out := renderMarkdown(t, "#### Deep Section")
	if !strings.Contains(out, "Deep Section") {
		t.Fatalf("expected H4 output to contain 'Deep Section', got: %q", out)
	}
}

func TestRenderHorizontalRule(t *testing.T) {
	out := renderMarkdown(t, "Above\n\n---\n\nBelow")
	if !strings.Contains(out, "─") {
		t.Fatalf("expected horizontal rule (─) in output, got: %q", out)
	}
	if !strings.Contains(out, "Above") || !strings.Contains(out, "Below") {
		t.Fatalf("expected text around horizontal rule, got: %q", out)
	}
}
