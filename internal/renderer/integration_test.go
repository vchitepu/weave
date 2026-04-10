package renderer

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func TestFullDocumentRender(t *testing.T) {
	input, err := os.ReadFile("../../testdata/full.md")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	th := theme.DarkTheme()
	r := New(th, 80)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(input, &buf)
	if err != nil {
		t.Fatalf("failed to render full document: %v", err)
	}

	out := buf.String()

	checks := []struct {
		desc string
		want string
	}{
		{"H1 heading", "Shine Test Document"},
		{"H4 heading", "Detailed Notes"},
		{"H5 heading", "Implementation Details"},
		{"H6 heading", "Footnotes"},
		{"bold text", "bold text"},
		{"italic text", "italic text"},
		{"inline code", "inline code"},
		{"link text", "Go"},
		{"link URL", "https://golang.org"},
		{"image alt", "screenshot"},
		{"code block content", "Println"},
		{"code block border", "╭"},
		{"language badge go", "go"},
		{"shell container", "go build"},
		{"tree content", "main.go"},
		{"diagram content", "Parser"},
		{"mermaid note", "diagram not rendered"},
		{"unordered bullet", "•"},
		{"ordered list", "1."},
		{"blockquote bar", "│"},
		{"table header", "Feature"},
		{"table content", "Headings"},
		{"table border", "─"},
		{"horizontal rule", "─"},
	}

	for _, c := range checks {
		if !strings.Contains(out, c.want) {
			t.Errorf("[%s] expected output to contain %q", c.desc, c.want)
		}
	}

	if len(out) < 500 {
		t.Errorf("expected substantial output, got only %d bytes", len(out))
	}
}

func TestRenderWithLightTheme(t *testing.T) {
	input, err := os.ReadFile("../../testdata/full.md")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	th := theme.LightTheme()
	r := New(th, 120)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(input, &buf)
	if err != nil {
		t.Fatalf("failed to render with light theme: %v", err)
	}

	if buf.Len() < 500 {
		t.Errorf("expected substantial output with light theme, got %d bytes", buf.Len())
	}
}

func TestRenderNarrowWidth(t *testing.T) {
	input, err := os.ReadFile("../../testdata/full.md")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	th := theme.DarkTheme()
	r := New(th, 40)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(input, &buf)
	if err != nil {
		t.Fatalf("failed to render at narrow width: %v", err)
	}

	if buf.Len() < 100 {
		t.Errorf("expected output at narrow width, got %d bytes", buf.Len())
	}
}

func TestHeadingLevelsProduceDistinctStyledOutput(t *testing.T) {
	input := []byte("### Same Text\n#### Same Text\n##### Same Text\n###### Same Text\n")

	th := theme.DarkTheme()
	r := New(th, 80)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)

	var buf bytes.Buffer
	err := md.Convert(input, &buf)
	if err != nil {
		t.Fatalf("failed to render heading levels: %v", err)
	}

	lines := strings.Split(buf.String(), "\n")
	matches := make([]string, 0, 4)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.Contains(line, "Same Text") {
			matches = append(matches, line)
		}
	}

	if len(matches) < 4 {
		t.Fatalf("expected at least 4 rendered heading lines containing %q, got %d", "Same Text", len(matches))
	}

	if matches[0] == matches[1] {
		t.Errorf("expected H3 and H4 rendered lines to be distinct, got %q", matches[0])
	}
	if matches[1] == matches[2] {
		t.Errorf("expected H4 and H5 rendered lines to be distinct, got %q", matches[1])
	}
	if matches[2] == matches[3] {
		t.Errorf("expected H5 and H6 rendered lines to be distinct, got %q", matches[2])
	}
	if matches[0] == matches[3] {
		t.Errorf("expected H3 and H6 rendered lines to be distinct, got %q", matches[0])
	}
}
