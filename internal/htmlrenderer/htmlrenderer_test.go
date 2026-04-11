package htmlrenderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// helper: render markdown string to HTML output string
func renderHTML(t *testing.T, input string) string {
	t.Helper()
	th := theme.DarkTheme()
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
	err := md.Convert([]byte(input), &buf)
	if err != nil {
		t.Fatalf("failed to convert markdown: %v", err)
	}
	return buf.String()
}

func TestHeading_H1(t *testing.T) {
	out := renderHTML(t, "# Hello World")
	if !strings.Contains(out, "<h1>") {
		t.Errorf("expected <h1> tag, got: %q", out)
	}
	if !strings.Contains(out, "Hello World") {
		t.Errorf("expected heading text, got: %q", out)
	}
	if !strings.Contains(out, "</h1>") {
		t.Errorf("expected </h1> closing tag, got: %q", out)
	}
}

func TestHeading_H2(t *testing.T) {
	out := renderHTML(t, "## Features")
	if !strings.Contains(out, "<h2>") {
		t.Errorf("expected <h2> tag, got: %q", out)
	}
	if !strings.Contains(out, "Features") {
		t.Errorf("expected heading text, got: %q", out)
	}
	if !strings.Contains(out, "</h2>") {
		t.Errorf("expected </h2> closing tag, got: %q", out)
	}
}

func TestHeading_H3Through6(t *testing.T) {
	tests := []struct {
		md  string
		tag string
	}{
		{"### H3 Heading", "h3"},
		{"#### H4 Heading", "h4"},
		{"##### H5 Heading", "h5"},
		{"###### H6 Heading", "h6"},
	}
	for _, tt := range tests {
		out := renderHTML(t, tt.md)
		if !strings.Contains(out, "<"+tt.tag+">") {
			t.Errorf("expected <%s> tag for %q, got: %q", tt.tag, tt.md, out)
		}
		if !strings.Contains(out, "</"+tt.tag+">") {
			t.Errorf("expected </%s> closing tag for %q, got: %q", tt.tag, tt.md, out)
		}
	}
}

func TestThematicBreak(t *testing.T) {
	out := renderHTML(t, "---")
	if !strings.Contains(out, "<hr>") {
		t.Errorf("expected <hr> tag, got: %q", out)
	}
}

func TestEmphasis(t *testing.T) {
	out := renderHTML(t, "*italic text*")
	if !strings.Contains(out, "<em>") {
		t.Errorf("expected <em> tag, got: %q", out)
	}
	if !strings.Contains(out, "italic text") {
		t.Errorf("expected emphasis text, got: %q", out)
	}
	if !strings.Contains(out, "</em>") {
		t.Errorf("expected </em> closing tag, got: %q", out)
	}
}

func TestStrong(t *testing.T) {
	out := renderHTML(t, "**bold text**")
	if !strings.Contains(out, "<strong>") {
		t.Errorf("expected <strong> tag, got: %q", out)
	}
	if !strings.Contains(out, "bold text") {
		t.Errorf("expected strong text, got: %q", out)
	}
	if !strings.Contains(out, "</strong>") {
		t.Errorf("expected </strong> closing tag, got: %q", out)
	}
}

func TestInlineCode(t *testing.T) {
	out := renderHTML(t, "Use `fmt.Println` here")
	if !strings.Contains(out, `class="inline-code"`) {
		t.Errorf("expected inline-code class, got: %q", out)
	}
	if !strings.Contains(out, "fmt.Println") {
		t.Errorf("expected code text, got: %q", out)
	}
}

func TestLink(t *testing.T) {
	out := renderHTML(t, "[Go](https://golang.org)")
	if !strings.Contains(out, `<a href="https://golang.org">`) {
		t.Errorf("expected link tag with href, got: %q", out)
	}
	if !strings.Contains(out, "Go</a>") {
		t.Errorf("expected link text, got: %q", out)
	}
}

func TestImage(t *testing.T) {
	out := renderHTML(t, "![screenshot](image.png)")
	if !strings.Contains(out, `<img src="image.png"`) {
		t.Errorf("expected img tag with src, got: %q", out)
	}
	if !strings.Contains(out, `alt="screenshot"`) {
		t.Errorf("expected alt attribute, got: %q", out)
	}
}

func TestStrikethrough(t *testing.T) {
	out := renderHTML(t, "~~deleted~~")
	if !strings.Contains(out, "<del>") {
		t.Errorf("expected <del> tag, got: %q", out)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected strikethrough text, got: %q", out)
	}
	if !strings.Contains(out, "</del>") {
		t.Errorf("expected </del> closing tag, got: %q", out)
	}
}

func TestFencedCodeBlock_WithLanguage(t *testing.T) {
	input := "```go\nfmt.Println(\"hello\")\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="code-container"`) {
		t.Errorf("expected code-container class, got: %q", out)
	}
	if !strings.Contains(out, `class="lang-badge"`) {
		t.Errorf("expected lang-badge, got: %q", out)
	}
	if !strings.Contains(out, "go</span>") {
		t.Errorf("expected go language badge, got: %q", out)
	}
	if !strings.Contains(out, "<pre>") {
		t.Errorf("expected <pre> tag, got: %q", out)
	}
	// chroma inline styles should be present for Go syntax
	if !strings.Contains(out, "Println") {
		t.Errorf("expected code content, got: %q", out)
	}
}

func TestFencedCodeBlock_NoLanguage(t *testing.T) {
	input := "```\nplain code\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="code-container"`) {
		t.Errorf("expected code-container class, got: %q", out)
	}
	if !strings.Contains(out, "plain code") {
		t.Errorf("expected code content, got: %q", out)
	}
	// No badge for no language
	if strings.Contains(out, `class="lang-badge"`) {
		t.Errorf("did not expect lang-badge for no-language code block, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerTree(t *testing.T) {
	input := "```tree\nfoo/\n  bar.go\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-tree") {
		t.Errorf("expected container-tree class, got: %q", out)
	}
	if !strings.Contains(out, ">tree</span>") {
		t.Errorf("expected tree badge, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerShell(t *testing.T) {
	input := "```bash\n$ go build ./cmd/weave\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-shell") {
		t.Errorf("expected container-shell class, got: %q", out)
	}
	if !strings.Contains(out, ">$</span>") {
		t.Errorf("expected $ badge, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerDiagram(t *testing.T) {
	input := "```diagram\n+---+\n| A |\n+---+\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-diagram") {
		t.Errorf("expected container-diagram class, got: %q", out)
	}
	if !strings.Contains(out, ">diagram</span>") {
		t.Errorf("expected diagram badge, got: %q", out)
	}
}

func TestBlockquote(t *testing.T) {
	out := renderHTML(t, "> quoted text")
	if !strings.Contains(out, "<blockquote>") {
		t.Errorf("expected <blockquote> tag, got: %q", out)
	}
	if !strings.Contains(out, "quoted text") {
		t.Errorf("expected blockquote text, got: %q", out)
	}
	if !strings.Contains(out, "</blockquote>") {
		t.Errorf("expected </blockquote> closing tag, got: %q", out)
	}
}

func TestUnorderedList(t *testing.T) {
	input := "- Apple\n- Banana\n- Cherry"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<ul>") {
		t.Errorf("expected <ul> tag, got: %q", out)
	}
	if !strings.Contains(out, "<li>") {
		t.Errorf("expected <li> tag, got: %q", out)
	}
	if !strings.Contains(out, "Apple") {
		t.Errorf("expected list item text, got: %q", out)
	}
	if !strings.Contains(out, "</ul>") {
		t.Errorf("expected </ul> closing tag, got: %q", out)
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. First\n2. Second\n3. Third"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<ol>") {
		t.Errorf("expected <ol> tag, got: %q", out)
	}
	if !strings.Contains(out, "<li>") {
		t.Errorf("expected <li> tag, got: %q", out)
	}
	if !strings.Contains(out, "First") {
		t.Errorf("expected list item text, got: %q", out)
	}
	if !strings.Contains(out, "</ol>") {
		t.Errorf("expected </ol> closing tag, got: %q", out)
	}
}

func TestTaskList_Checked(t *testing.T) {
	input := "- [x] Done task"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="task-item"`) {
		t.Errorf("expected task-item class, got: %q", out)
	}
	if !strings.Contains(out, `class="checkbox checked"`) {
		t.Errorf("expected checked checkbox, got: %q", out)
	}
	if !strings.Contains(out, "&#x2713;") {
		t.Errorf("expected checkmark entity, got: %q", out)
	}
}

func TestTaskList_Unchecked(t *testing.T) {
	input := "- [ ] Pending task"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="task-item"`) {
		t.Errorf("expected task-item class, got: %q", out)
	}
	if !strings.Contains(out, `class="checkbox unchecked"`) {
		t.Errorf("expected unchecked checkbox, got: %q", out)
	}
	if !strings.Contains(out, "&#x25CB;") {
		t.Errorf("expected circle entity, got: %q", out)
	}
}

func TestTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<table>") {
		t.Errorf("expected <table> tag, got: %q", out)
	}
	if !strings.Contains(out, "<thead>") {
		t.Errorf("expected <thead> tag, got: %q", out)
	}
	if !strings.Contains(out, "<th>") {
		t.Errorf("expected <th> tag, got: %q", out)
	}
	if !strings.Contains(out, "Name") {
		t.Errorf("expected header text, got: %q", out)
	}
	if !strings.Contains(out, "<tbody>") {
		t.Errorf("expected <tbody> tag, got: %q", out)
	}
	if !strings.Contains(out, "<td>") {
		t.Errorf("expected <td> tag, got: %q", out)
	}
	if !strings.Contains(out, "Alice") {
		t.Errorf("expected cell text, got: %q", out)
	}
	if !strings.Contains(out, "</table>") {
		t.Errorf("expected </table> closing tag, got: %q", out)
	}
}
