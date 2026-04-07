package renderer

import (
	"strings"
	"testing"
)

func TestRenderBold(t *testing.T) {
	out := renderMarkdown(t, "This is **bold** text")
	if !strings.Contains(out, "bold") {
		t.Fatalf("expected bold text in output, got: %q", out)
	}
}

func TestRenderItalic(t *testing.T) {
	out := renderMarkdown(t, "This is *italic* text")
	if !strings.Contains(out, "italic") {
		t.Fatalf("expected italic text in output, got: %q", out)
	}
}

func TestRenderInlineCode(t *testing.T) {
	out := renderMarkdown(t, "Use `fmt.Println` here")
	if !strings.Contains(out, "fmt.Println") {
		t.Fatalf("expected inline code in output, got: %q", out)
	}
}

func TestRenderLink(t *testing.T) {
	out := renderMarkdown(t, "[Go](https://golang.org)")
	if !strings.Contains(out, "Go") {
		t.Fatalf("expected link text in output, got: %q", out)
	}
	if !strings.Contains(out, "https://golang.org") {
		t.Fatalf("expected URL in output, got: %q", out)
	}
}

func TestRenderImage(t *testing.T) {
	out := renderMarkdown(t, "![screenshot](img.png)")
	if !strings.Contains(out, "screenshot") {
		t.Fatalf("expected image alt text in output, got: %q", out)
	}
}

func TestRenderStrikethrough(t *testing.T) {
	out := renderMarkdown(t, "This is ~~deleted~~ text")
	if !strings.Contains(out, "deleted") {
		t.Fatalf("expected strikethrough text in output, got: %q", out)
	}
}
