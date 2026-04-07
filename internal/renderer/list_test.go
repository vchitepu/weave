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

func TestRenderListItemsSeparated(t *testing.T) {
	input := "- Item A\n- Item B\n- Item C"
	out := renderMarkdown(t, input)
	// Each item should be on its own line
	if strings.Contains(out, "A•") || strings.Contains(out, "A  •") {
		t.Fatalf("list items should not concatenate on same line, got: %q", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	bulletLines := 0
	for _, line := range lines {
		if strings.Contains(line, "•") {
			bulletLines++
		}
	}
	if bulletLines < 3 {
		t.Fatalf("expected at least 3 lines with bullets, got %d in output: %q", bulletLines, out)
	}
}
