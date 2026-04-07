package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vinaychitepu/shine/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// helper that enables the table extension
func renderMarkdownWithTables(t *testing.T, input string) string {
	t.Helper()
	th := theme.DarkTheme()
	r := New(th, 80)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, 100),
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

func TestRenderSimpleTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	out := renderMarkdownWithTables(t, input)
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
	out := renderMarkdownWithTables(t, input)
	if !strings.Contains(out, "Left") {
		t.Fatalf("expected 'Left' in table output, got: %q", out)
	}
}
