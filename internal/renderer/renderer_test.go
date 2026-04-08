package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vinaychitepu/weave/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// helper: render markdown string to terminal output string
func renderMarkdown(t *testing.T, input string) string {
	t.Helper()
	th := theme.DarkTheme()
	r := New(th, 80)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)
	var buf bytes.Buffer
	err := md.Convert([]byte(input), &buf)
	if err != nil {
		t.Fatalf("failed to convert markdown: %v", err)
	}
	return buf.String()
}

func TestRenderPlainText(t *testing.T) {
	out := renderMarkdown(t, "Hello world")
	if !strings.Contains(out, "Hello world") {
		t.Fatalf("expected output to contain 'Hello world', got: %q", out)
	}
}

func TestRenderMultipleParagraphs(t *testing.T) {
	out := renderMarkdown(t, "First paragraph.\n\nSecond paragraph.")
	if !strings.Contains(out, "First paragraph.") {
		t.Fatalf("expected output to contain 'First paragraph.', got: %q", out)
	}
	if !strings.Contains(out, "Second paragraph.") {
		t.Fatalf("expected output to contain 'Second paragraph.', got: %q", out)
	}
}
