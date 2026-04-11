package htmlrenderer

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

func renderWithTheme(t *testing.T, input []byte, th theme.Theme) string {
	t.Helper()
	r := New(th)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
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
		t.Fatalf("failed to convert markdown: %v", err)
	}
	return buf.String()
}

func TestIntegration_FullDocument(t *testing.T) {
	input, err := os.ReadFile("../../testdata/full.md")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	th := theme.DarkTheme()
	out := renderWithTheme(t, input, th)

	checks := []struct {
		desc string
		want string
	}{
		{"H1 heading", "<h1>"},
		{"H1 text", "Shine Test Document"},
		{"H2 heading", "<h2>"},
		{"H3 heading", "<h3>"},
		{"H4 heading", "<h4>"},
		{"H5 heading", "<h5>"},
		{"H6 heading", "<h6>"},
		{"bold text", "<strong>"},
		{"italic text", "<em>"},
		{"inline code", `class="inline-code"`},
		{"link tag", "<a href="},
		{"link url", "https://golang.org"},
		{"image tag", "<img src="},
		{"code container", `class="code-container`},
		{"language badge", `class="lang-badge"`},
		{"pre tag", "<pre>"},
		{"shell container", "container-shell"},
		{"tree container", "container-tree"},
		{"diagram container", "container-diagram"},
		{"blockquote", "<blockquote>"},
		{"unordered list", "<ul>"},
		{"ordered list", "<ol>"},
		{"list item", "<li>"},
		{"table", "<table>"},
		{"table header", "<thead>"},
		{"table header cell", "<th>"},
		{"table body", "<tbody>"},
		{"table cell", "<td>"},
		{"horizontal rule", "<hr>"},
		{"task item", `class="task-item"`},
		{"checked task", `class="checkbox checked"`},
		{"unchecked task", `class="checkbox unchecked"`},
		{"code content", "Println"},
		{"table content", "Headings"},
		{"list content", "First item"},
		{"blockquote content", "blockquote"},
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

func TestIntegration_DarkAndLightThemesDiffer(t *testing.T) {
	// Use a code block with a language so chroma produces different inline styles
	input := []byte("```go\npackage main\n\nfunc main() {}\n```\n")

	darkOut := renderWithTheme(t, input, theme.DarkTheme())
	lightOut := renderWithTheme(t, input, theme.LightTheme())

	if darkOut == lightOut {
		t.Errorf("expected dark and light theme outputs to differ for syntax-highlighted code")
	}
}
