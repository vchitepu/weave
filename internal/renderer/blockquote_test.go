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
