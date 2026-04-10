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

func TestRenderH4(t *testing.T) {
	out := renderMarkdown(t, "#### Deep Section")
	if !strings.Contains(out, "Deep Section") {
		t.Fatalf("expected H4 output to contain 'Deep Section', got: %q", out)
	}
}

func TestRenderH5(t *testing.T) {
	out := renderMarkdown(t, "##### Deep Section")
	if !strings.Contains(out, "Deep Section") {
		t.Fatalf("expected H5 output to contain 'Deep Section', got: %q", out)
	}
}

func TestRenderH6(t *testing.T) {
	out := renderMarkdown(t, "###### Deep Section")
	if !strings.Contains(out, "Deep Section") {
		t.Fatalf("expected H6 output to contain 'Deep Section', got: %q", out)
	}
}

func TestRenderH4H5H6ProduceDistinctOutput(t *testing.T) {
	h3 := renderMarkdown(t, "### Same Text")
	h4 := renderMarkdown(t, "#### Same Text")
	h5 := renderMarkdown(t, "##### Same Text")
	h6 := renderMarkdown(t, "###### Same Text")

	if h4 == h3 {
		t.Fatalf("expected H4 output to differ from H3 output, both were: %q", h4)
	}
	if h5 == h3 {
		t.Fatalf("expected H5 output to differ from H3 output, both were: %q", h5)
	}
	if h6 == h3 {
		t.Fatalf("expected H6 output to differ from H3 output, both were: %q", h6)
	}
	if h4 == h5 {
		t.Fatalf("expected H4 output to differ from H5 output, both were: %q", h4)
	}
	if h5 == h6 {
		t.Fatalf("expected H5 output to differ from H6 output, both were: %q", h5)
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

func TestRenderH1WithInlineFormatting(t *testing.T) {
	out := renderMarkdown(t, "# Hello **world**")
	if !strings.Contains(out, "Hello") {
		t.Fatalf("expected 'Hello' in heading output, got: %q", out)
	}
	if !strings.Contains(out, "world") {
		t.Fatalf("expected 'world' from bold text in heading output, got: %q", out)
	}
}
