package renderer

import (
	"strings"
	"testing"
)

func TestRenderUnorderedList(t *testing.T) {
	input := "- Item one\n- Item two\n- Item three"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "•") {
		t.Fatalf("expected bullet '•' in output, got: %q", out)
	}
	if !strings.Contains(out, "Item one") {
		t.Fatalf("expected 'Item one' in output, got: %q", out)
	}
	if !strings.Contains(out, "Item three") {
		t.Fatalf("expected 'Item three' in output, got: %q", out)
	}
}

func TestRenderOrderedList(t *testing.T) {
	input := "1. First\n2. Second\n3. Third"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "1.") {
		t.Fatalf("expected '1.' in output, got: %q", out)
	}
	if !strings.Contains(out, "First") {
		t.Fatalf("expected 'First' in output, got: %q", out)
	}
}

func TestRenderNestedList(t *testing.T) {
	input := "- Outer\n  - Inner"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "Outer") {
		t.Fatalf("expected 'Outer' in output, got: %q", out)
	}
	if !strings.Contains(out, "Inner") {
		t.Fatalf("expected 'Inner' in output, got: %q", out)
	}
}
