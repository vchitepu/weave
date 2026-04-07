package renderer

import (
	"strings"
	"testing"
)

func TestRenderSimpleTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	out := renderMarkdown(t, input)
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
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "Left") {
		t.Fatalf("expected 'Left' in table output, got: %q", out)
	}
}
