package renderer

import (
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
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
	th := theme.DarkTheme()
	th.H3 = th.H3.Transform(func(s string) string { return "H3:" + s })
	th.H4 = th.H4.Transform(func(s string) string { return "H4:" + s })
	th.H5 = th.H5.Transform(func(s string) string { return "H5:" + s })
	th.H6 = th.H6.Transform(func(s string) string { return "H6:" + s })

	h3 := renderMarkdownWithTheme(t, "### Same Text", th)
	h4 := renderMarkdownWithTheme(t, "#### Same Text", th)
	h5 := renderMarkdownWithTheme(t, "##### Same Text", th)
	h6 := renderMarkdownWithTheme(t, "###### Same Text", th)

	if !strings.Contains(h3, "H3:Same Text") {
		t.Fatalf("expected H3 marker in output, got: %q", h3)
	}
	if !strings.Contains(h4, "H4:Same Text") {
		t.Fatalf("expected H4 marker in output, got: %q", h4)
	}
	if !strings.Contains(h5, "H5:Same Text") {
		t.Fatalf("expected H5 marker in output, got: %q", h5)
	}
	if !strings.Contains(h6, "H6:Same Text") {
		t.Fatalf("expected H6 marker in output, got: %q", h6)
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
