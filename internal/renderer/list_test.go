package renderer

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
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

func TestRenderLongListItemWrapsWithinRendererWidth(t *testing.T) {
	input := "- OS updates: Standard Ubuntu apt security and package updates. OpenLM system components (orchestrator, event bus, LiteLLM) update via the OpenLM apt repo."
	out := renderMarkdown(t, input)

	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if w := lipgloss.Width(line); w > 80 {
			t.Fatalf("expected list output line width <= 80, got %d in line %q", w, line)
		}
	}
}

func TestRenderTaskListChecked(t *testing.T) {
	input := "- [x] Done item"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "✓") {
		t.Fatalf("expected checked symbol '✓' in output, got: %q", out)
	}
	if !strings.Contains(out, "Done item") {
		t.Fatalf("expected 'Done item' in output, got: %q", out)
	}
	// Should NOT contain a bullet — the checkbox replaces it
	if strings.Contains(out, "•") {
		t.Fatalf("task list item should not have bullet '•', got: %q", out)
	}
}

func TestRenderTaskListUnchecked(t *testing.T) {
	input := "- [ ] Pending item"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "○") {
		t.Fatalf("expected unchecked symbol '○' in output, got: %q", out)
	}
	if !strings.Contains(out, "Pending item") {
		t.Fatalf("expected 'Pending item' in output, got: %q", out)
	}
	if strings.Contains(out, "•") {
		t.Fatalf("task list item should not have bullet '•', got: %q", out)
	}
}

func TestRenderTaskListMixed(t *testing.T) {
	input := "- [x] Done\n- [ ] Todo\n- [X] Also done"
	out := renderMarkdown(t, input)
	if !strings.Contains(out, "✓") {
		t.Fatalf("expected '✓' in output, got: %q", out)
	}
	if !strings.Contains(out, "○") {
		t.Fatalf("expected '○' in output, got: %q", out)
	}
	if !strings.Contains(out, "Done") {
		t.Fatalf("expected 'Done' in output, got: %q", out)
	}
	if !strings.Contains(out, "Todo") {
		t.Fatalf("expected 'Todo' in output, got: %q", out)
	}
}

func TestRenderLongTaskListItemWrapsWithinRendererWidth(t *testing.T) {
	input := "- [x] This is a very long task list item that should wrap within the renderer width boundary without exceeding eighty columns of terminal output"
	out := renderMarkdown(t, input)

	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if w := lipgloss.Width(line); w > 80 {
			t.Fatalf("expected task list output line width <= 80, got %d in line %q", w, line)
		}
	}
}

func TestRenderNestedListItemsAreOnSeparateLines(t *testing.T) {
	input := "- First item\n    - Second item\n        - Nested item\n    - Another nested"
	out := renderMarkdown(t, input)

	if strings.Contains(out, "Second itemNested item") {
		t.Fatalf("expected nested list items to render on separate lines, got: %q", out)
	}
	if !strings.Contains(out, "Nested item") {
		t.Fatalf("expected nested item text in output, got: %q", out)
	}
}
