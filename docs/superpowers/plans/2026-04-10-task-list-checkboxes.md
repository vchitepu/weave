# Task List Checkboxes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Render GFM task list checkboxes (`- [x]` / `- [ ]`) as styled terminal symbols.

**Architecture:** Enable goldmark's `extension.TaskList` parser to produce `ast.TaskCheckBox` inline nodes. Register a custom `renderTaskCheckBox` handler in the renderer that writes themed `✓` or `○` symbols. Add `TaskChecked`/`TaskUnchecked` styles to the theme struct.

**Tech Stack:** Go, goldmark, lipgloss, chroma (existing)

---

### Task 1: Add theme styles for task checkboxes

**Files:**
- Modify: `internal/theme/theme.go:6-49` (Theme struct), `internal/theme/theme.go:52-89` (DarkTheme), `internal/theme/theme.go:92-129` (LightTheme)
- Test: `internal/theme/theme_test.go`

- [ ] **Step 1: Write the failing test**

Add to `internal/theme/theme_test.go`:

```go
func TestDarkThemeTaskCheckboxStyles(t *testing.T) {
	th := DarkTheme()
	if th.TaskChecked.GetForeground() == nil {
		t.Fatal("DarkTheme TaskChecked should have a foreground color")
	}
	if got := colorString(th.TaskChecked.GetForeground()); got != "#6BBF8A" {
		t.Fatalf("DarkTheme TaskChecked color = %q, want %q", got, "#6BBF8A")
	}
	if th.TaskUnchecked.GetForeground() == nil {
		t.Fatal("DarkTheme TaskUnchecked should have a foreground color")
	}
	if got := colorString(th.TaskUnchecked.GetForeground()); got != "#5A5F6B" {
		t.Fatalf("DarkTheme TaskUnchecked color = %q, want %q", got, "#5A5F6B")
	}
}

func TestLightThemeTaskCheckboxStyles(t *testing.T) {
	th := LightTheme()
	if got := colorString(th.TaskChecked.GetForeground()); got != "#3A8C5C" {
		t.Fatalf("LightTheme TaskChecked color = %q, want %q", got, "#3A8C5C")
	}
	if got := colorString(th.TaskUnchecked.GetForeground()); got != "#9BA3B2" {
		t.Fatalf("LightTheme TaskUnchecked color = %q, want %q", got, "#9BA3B2")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/theme/ -run TestDarkThemeTaskCheckboxStyles -v`
Expected: FAIL — `TaskChecked` field does not exist

- [ ] **Step 3: Add fields to Theme struct**

In `internal/theme/theme.go`, add two new fields to the `Theme` struct after the `HorizontalRule` field:

```go
	// Task list checkboxes
	TaskChecked   lipgloss.Style
	TaskUnchecked lipgloss.Style
```

- [ ] **Step 4: Wire DarkTheme values**

In `internal/theme/theme.go` `DarkTheme()`, add before the closing brace:

```go
		TaskChecked:   lipgloss.NewStyle().Foreground(lipgloss.Color("#6BBF8A")),
		TaskUnchecked: lipgloss.NewStyle().Foreground(lipgloss.Color("#5A5F6B")),
```

- [ ] **Step 5: Wire LightTheme values**

In `internal/theme/theme.go` `LightTheme()`, add before the closing brace:

```go
		TaskChecked:   lipgloss.NewStyle().Foreground(lipgloss.Color("#3A8C5C")),
		TaskUnchecked: lipgloss.NewStyle().Foreground(lipgloss.Color("#9BA3B2")),
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/theme/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add internal/theme/theme.go internal/theme/theme_test.go
git commit -m "feat: add TaskChecked and TaskUnchecked styles to theme"
```

---

### Task 2: Enable TaskList extension in the parser

**Files:**
- Modify: `cmd/weave/main.go:88`
- Modify: `internal/renderer/renderer_test.go:26` (test helper)
- Modify: `internal/renderer/integration_test.go:25,93,123` (integration goldmark setup)

- [ ] **Step 1: Add extension.TaskList to main.go**

In `cmd/weave/main.go`, change line 88 from:

```go
		goldmark.WithExtensions(extension.Table, extension.Strikethrough),
```

to:

```go
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
```

- [ ] **Step 2: Add extension.TaskList to the test helper**

In `internal/renderer/renderer_test.go`, change line 26 from:

```go
		goldmark.WithExtensions(extension.Table, extension.Strikethrough),
```

to:

```go
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
```

- [ ] **Step 3: Add extension.TaskList to integration test goldmark setups**

In `internal/renderer/integration_test.go`, there are four goldmark.New calls (lines 25, 93, 123, 154). Each has `goldmark.WithExtensions(extension.Table)`. Change all four from:

```go
		goldmark.WithExtensions(extension.Table),
```

to:

```go
		goldmark.WithExtensions(extension.Table, extension.TaskList),
```

- [ ] **Step 4: Run existing tests to confirm nothing breaks**

Run: `go test ./... -v`
Expected: All PASS (goldmark now parses task lists but our renderer doesn't handle the node yet — it will be silently ignored since unregistered nodes are skipped)

- [ ] **Step 5: Commit**

```bash
git add cmd/weave/main.go internal/renderer/renderer_test.go internal/renderer/integration_test.go
git commit -m "feat: enable goldmark TaskList extension in parser"
```

---

### Task 3: Implement renderTaskCheckBox and bullet suppression

**Files:**
- Modify: `internal/renderer/list.go` (add renderTaskCheckBox, modify renderListItem)
- Modify: `internal/renderer/renderer.go:96-124` (register handler in RegisterFuncs)
- Test: `internal/renderer/list_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `internal/renderer/list_test.go`:

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/renderer/ -run TestRenderTaskList -v`
Expected: FAIL — `✓` not found in output (checkbox node is silently skipped, bullet is still rendered)

- [ ] **Step 3: Register the handler in RegisterFuncs**

In `internal/renderer/renderer.go`, add an import for the extension AST package (it is already imported as `east`). Then in the `RegisterFuncs` method, after the strikethrough registration (line 115), add:

```go
	// Task list checkboxes (from goldmark extension)
	reg.Register(east.KindTaskCheckBox, r.renderTaskCheckBox)
```

- [ ] **Step 4: Implement renderTaskCheckBox in list.go**

In `internal/renderer/list.go`, add the import for the goldmark extension AST package at the top:

```go
	east "github.com/yuin/goldmark/extension/ast"
```

Then add the handler function at the end of the file:

```go
func (r *Renderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*east.TaskCheckBox)
	if n.IsChecked {
		_, _ = w.WriteString(r.theme.TaskChecked.Render("✓") + " ")
	} else {
		_, _ = w.WriteString(r.theme.TaskUnchecked.Render("○") + " ")
	}
	return ast.WalkContinue, nil
}
```

- [ ] **Step 5: Add bullet suppression in renderListItem**

In `internal/renderer/list.go`, modify the `renderListItem` function. After line 62 (`list := node.Parent().(*ast.List)`), add task checkbox detection. The complete bullet/number section (lines 62-79) becomes:

```go
	// Determine bullet or number
	list := node.Parent().(*ast.List)

	// Check if this list item has a task checkbox — if so, suppress the bullet
	// and write only the indent prefix. The checkbox symbol is rendered by
	// renderTaskCheckBox.
	isTask := false
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		// TaskCheckBox is an inline node inside a TextBlock/Paragraph child
		for inline := child.FirstChild(); inline != nil; inline = inline.NextSibling() {
			if inline.Kind() == east.KindTaskCheckBox {
				isTask = true
				break
			}
		}
		if isTask {
			break
		}
	}

	if isTask {
		prefix := fmt.Sprintf("%s%s", pad, indent)
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix)+2)
		_, _ = w.WriteString(prefix)
	} else if list.IsOrdered() {
		pos := 1
		for sib := node.PreviousSibling(); sib != nil; sib = sib.PreviousSibling() {
			pos++
		}
		start := list.Start
		if start > 0 {
			pos = start + pos - 1
		}
		prefix := fmt.Sprintf("%s%s%d. ", pad, indent, pos)
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix))
		_, _ = w.WriteString(prefix)
	} else {
		prefix := fmt.Sprintf("%s%s• ", pad, indent)
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix))
		_, _ = w.WriteString(prefix)
	}
```

Note: the `+2` in `lipgloss.Width(prefix)+2` accounts for the checkbox symbol (`✓ ` or `○ `) that `renderTaskCheckBox` will write, so wrapped continuation lines align correctly.

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/renderer/ -run TestRenderTaskList -v`
Expected: All PASS

- [ ] **Step 7: Run full test suite**

Run: `go test ./... -v`
Expected: All PASS

- [ ] **Step 8: Commit**

```bash
git add internal/renderer/list.go internal/renderer/renderer.go
git commit -m "feat: render task list checkboxes with themed symbols"
```

---

### Task 4: Add task list items to test fixture and integration test

**Files:**
- Modify: `testdata/full.md`
- Modify: `internal/renderer/integration_test.go`

- [ ] **Step 1: Add task list section to fixture**

In `testdata/full.md`, add a new section before the final `---` line (before line 100):

```markdown
## Task Lists

- [x] Completed task
- [ ] Pending task
- [x] Another done item
```

- [ ] **Step 2: Add integration test checks**

In `internal/renderer/integration_test.go`, in the `TestFullDocumentRender` function, add to the `checks` slice:

```go
		{"task checked", "✓"},
		{"task unchecked", "○"},
		{"task list text", "Completed task"},
```

- [ ] **Step 3: Run integration tests**

Run: `go test ./internal/renderer/ -run TestFullDocumentRender -v`
Expected: PASS

- [ ] **Step 4: Run full test suite**

Run: `go test ./... -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add testdata/full.md internal/renderer/integration_test.go
git commit -m "test: add task list checkboxes to fixture and integration tests"
```

---

### Task 5: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add task list to feature list**

In `README.md`, add a line after `- Ordered and unordered lists with nested indentation` (line 12):

```markdown
- Task list checkboxes with checked/unchecked symbols
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add task list checkboxes to README features"
```
