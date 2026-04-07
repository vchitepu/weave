package renderer

import (
	"strings"
	"testing"
)

func TestRenderFencedCodeBlock(t *testing.T) {
	input := "```go\nfunc main() {}\n```"
	out := renderMarkdown(t, input)
	// Syntax highlighting inserts ANSI escape codes, so check for individual tokens
	if !strings.Contains(out, "func") {
		t.Fatalf("expected 'func' in output, got: %q", out)
	}
	if !strings.Contains(out, "main") {
		t.Fatalf("expected 'main' in output, got: %q", out)
	}
	if !strings.Contains(out, "╭") || !strings.Contains(out, "╰") {
		t.Fatalf("expected rounded border in output, got: %q", out)
	}
	if !strings.Contains(out, "go") {
		t.Fatalf("expected language badge 'go' in output, got: %q", out)
	}
}

func TestRenderCodeBlockNoLanguage(t *testing.T) {
	input := "```\nhello world\n```"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "hello world") {
		t.Fatalf("expected code in output, got: %q", out)
	}
	if !strings.Contains(out, "╭") {
		t.Fatalf("expected border in output, got: %q", out)
	}
}

func TestDetectContainerType(t *testing.T) {
	tests := []struct {
		lang     string
		expected containerType
	}{
		{"go", containerCode},
		{"python", containerCode},
		{"tree", containerTree},
		{"ascii", containerDiagram},
		{"diagram", containerDiagram},
		{"art", containerDiagram},
		{"mermaid", containerDiagram},
		{"bash", containerShell},
		{"sh", containerShell},
		{"shell", containerShell},
		{"console", containerShell},
		{"terminal", containerShell},
		{"", containerCode},
	}
	for _, tt := range tests {
		got := detectContainer(tt.lang)
		if got != tt.expected {
			t.Errorf("detectContainer(%q) = %v, want %v", tt.lang, got, tt.expected)
		}
	}
}

func TestRenderTreeContainer(t *testing.T) {
	input := "```tree\nsrc/\n├── main.go\n└── lib/\n```"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "tree") {
		t.Fatalf("expected 'tree' label in output, got: %q", out)
	}
	if !strings.Contains(out, "main.go") {
		t.Fatalf("expected tree content in output, got: %q", out)
	}
}

func TestRenderShellContainer(t *testing.T) {
	input := "```bash\n$ go build ./...\n```"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "go build") {
		t.Fatalf("expected shell content in output, got: %q", out)
	}
}

func TestRenderMermaidContainer(t *testing.T) {
	input := "```mermaid\ngraph TD\n    A --> B\n```"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "diagram not rendered") {
		t.Fatalf("expected 'diagram not rendered' note for mermaid, got: %q", out)
	}
	if !strings.Contains(out, "graph TD") {
		t.Fatalf("expected mermaid source in output, got: %q", out)
	}
}
