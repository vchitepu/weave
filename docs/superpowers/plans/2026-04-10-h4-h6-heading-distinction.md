# H4–H6 Heading Distinction Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add visually distinct H4, H5, and H6 heading styles using progressively dimmer colors derived from the H3 palette.

**Architecture:** Three new `lipgloss.Style` fields (`H4`, `H5`, `H6`) are added to the `Theme` struct, populated in both `DarkTheme()` and `LightTheme()`. The renderer's heading dispatch switches from a fallback default to explicit per-level cases.

**Tech Stack:** Go, lipgloss (terminal styling), goldmark (markdown AST), standard testing

---

### Task 1: Add H4/H5/H6 fields to Theme struct and theme constructors

**Files:**
- Modify: `internal/theme/theme.go:14-18` (struct fields)
- Modify: `internal/theme/theme.go:57-60` (DarkTheme H4/H5/H6)
- Modify: `internal/theme/theme.go:94-97` (LightTheme H4/H5/H6)
- Test: `internal/theme/theme_test.go`

- [ ] **Step 1: Write failing tests for H4/H5/H6 theme fields**

Add to `internal/theme/theme_test.go` after the existing `TestDarkThemeNotNil` function:

```go
func TestDarkThemeH4H5H6(t *testing.T) {
	th := DarkTheme()

	if th.H4.GetForeground() == nil {
		t.Fatal("DarkTheme H4 should have a foreground color")
	}
	if got := colorString(th.H4.GetForeground()); got != "#A88A55" {
		t.Fatalf("DarkTheme H4 color = %q, want %q", got, "#A88A55")
	}

	if th.H5.GetForeground() == nil {
		t.Fatal("DarkTheme H5 should have a foreground color")
	}
	if got := colorString(th.H5.GetForeground()); got != "#876C42" {
		t.Fatalf("DarkTheme H5 color = %q, want %q", got, "#876C42")
	}

	if th.H6.GetForeground() == nil {
		t.Fatal("DarkTheme H6 should have a foreground color")
	}
	if got := colorString(th.H6.GetForeground()); got != "#665030" {
		t.Fatalf("DarkTheme H6 color = %q, want %q", got, "#665030")
	}
}

func TestLightThemeH4H5H6(t *testing.T) {
	th := LightTheme()

	if th.H4.GetForeground() == nil {
		t.Fatal("LightTheme H4 should have a foreground color")
	}
	if got := colorString(th.H4.GetForeground()); got != "#A07850" {
		t.Fatalf("LightTheme H4 color = %q, want %q", got, "#A07850")
	}

	if th.H5.GetForeground() == nil {
		t.Fatal("LightTheme H5 should have a foreground color")
	}
	if got := colorString(th.H5.GetForeground()); got != "#B38A63" {
		t.Fatalf("LightTheme H5 color = %q, want %q", got, "#B38A63")
	}

	if th.H6.GetForeground() == nil {
		t.Fatal("LightTheme H6 should have a foreground color")
	}
	if got := colorString(th.H6.GetForeground()); got != "#C69C78" {
		t.Fatalf("LightTheme H6 color = %q, want %q", got, "#C69C78")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/theme/ -run 'TestDarkThemeH4H5H6|TestLightThemeH4H5H6' -v`
Expected: FAIL — `H4`, `H5`, `H6` fields don't exist on Theme yet.

- [ ] **Step 3: Add H4/H5/H6 fields to Theme struct**

In `internal/theme/theme.go`, add three fields after line 17 (`H3 lipgloss.Style`):

```go
	H4          lipgloss.Style
	H5          lipgloss.Style
	H6          lipgloss.Style
```

- [ ] **Step 4: Populate H4/H5/H6 in DarkTheme()**

In `internal/theme/theme.go` `DarkTheme()`, add after the `H3` line (line 59):

```go
		H4:          lipgloss.NewStyle().Foreground(lipgloss.Color("#A88A55")),
		H5:          lipgloss.NewStyle().Foreground(lipgloss.Color("#876C42")),
		H6:          lipgloss.NewStyle().Foreground(lipgloss.Color("#665030")),
```

- [ ] **Step 5: Populate H4/H5/H6 in LightTheme()**

In `internal/theme/theme.go` `LightTheme()`, add after the `H3` line (line 96):

```go
		H4:          lipgloss.NewStyle().Foreground(lipgloss.Color("#A07850")),
		H5:          lipgloss.NewStyle().Foreground(lipgloss.Color("#B38A63")),
		H6:          lipgloss.NewStyle().Foreground(lipgloss.Color("#C69C78")),
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/theme/ -v`
Expected: ALL PASS — including new `TestDarkThemeH4H5H6` and `TestLightThemeH4H5H6`.

- [ ] **Step 7: Commit**

```bash
git add internal/theme/theme.go internal/theme/theme_test.go
git commit -m "feat: add H4, H5, H6 style fields to Theme struct"
```

---

### Task 2: Update renderer heading dispatch to use H4/H5/H6

**Files:**
- Modify: `internal/renderer/headings.go:20-27` (switch statement)
- Test: `internal/renderer/headings_test.go`

- [ ] **Step 1: Update existing test and write new failing tests for H4/H5/H6 rendering**

In `internal/renderer/headings_test.go`, replace `TestRenderH4UsesH3Style` and add two new tests:

Replace the existing `TestRenderH4UsesH3Style` (lines 33-38) with:

```go
func TestRenderH4(t *testing.T) {
	out := renderMarkdown(t, "#### Deep Section")
	if !strings.Contains(out, "Deep Section") {
		t.Fatalf("expected H4 output to contain 'Deep Section', got: %q", out)
	}
}

func TestRenderH5(t *testing.T) {
	out := renderMarkdown(t, "##### Deeper Section")
	if !strings.Contains(out, "Deeper Section") {
		t.Fatalf("expected H5 output to contain 'Deeper Section', got: %q", out)
	}
}

func TestRenderH6(t *testing.T) {
	out := renderMarkdown(t, "###### Deepest Section")
	if !strings.Contains(out, "Deepest Section") {
		t.Fatalf("expected H6 output to contain 'Deepest Section', got: %q", out)
	}
}

func TestRenderH4H5H6ProduceDistinctOutput(t *testing.T) {
	h4 := renderMarkdown(t, "#### H4 Heading")
	h5 := renderMarkdown(t, "##### H5 Heading")
	h6 := renderMarkdown(t, "###### H6 Heading")

	if h4 == h5 {
		t.Fatal("H4 and H5 should produce distinct output")
	}
	if h5 == h6 {
		t.Fatal("H5 and H6 should produce distinct output")
	}
	if h4 == h6 {
		t.Fatal("H4 and H6 should produce distinct output")
	}
}
```

- [ ] **Step 2: Run tests to verify the distinction test fails**

Run: `go test ./internal/renderer/ -run 'TestRenderH4H5H6ProduceDistinctOutput' -v`
Expected: FAIL — all three currently use `r.theme.H3`, so H4/H5/H6 produce identical output (except text content; the distinct-output test compares full renders with different text so it would pass trivially). Let's verify the test logic: since each render has different text content ("H4 Heading" vs "H5 Heading"), the outputs will differ even with same styling. We need a better test.

Actually, a better approach: render the same text with different heading levels and compare.

Replace the `TestRenderH4H5H6ProduceDistinctOutput` test above with:

```go
func TestRenderH4H5H6ProduceDistinctOutput(t *testing.T) {
	h3 := renderMarkdown(t, "### Same Text")
	h4 := renderMarkdown(t, "#### Same Text")
	h5 := renderMarkdown(t, "##### Same Text")
	h6 := renderMarkdown(t, "###### Same Text")

	if h4 == h3 {
		t.Fatal("H4 and H3 should produce distinct output")
	}
	if h5 == h3 {
		t.Fatal("H5 and H3 should produce distinct output")
	}
	if h6 == h3 {
		t.Fatal("H6 and H3 should produce distinct output")
	}
	if h4 == h5 {
		t.Fatal("H4 and H5 should produce distinct output")
	}
	if h5 == h6 {
		t.Fatal("H5 and H6 should produce distinct output")
	}
}
```

- [ ] **Step 3: Run tests to verify the distinction test fails**

Run: `go test ./internal/renderer/ -run 'TestRenderH4H5H6ProduceDistinctOutput' -v`
Expected: FAIL — `H4 and H3 should produce distinct output` because all still use `r.theme.H3`.

- [ ] **Step 4: Update the heading switch in headings.go**

In `internal/renderer/headings.go`, replace lines 20-27:

```go
	var styled string
	switch {
	case n.Level == 1:
		styled = r.theme.H1.Render(text)
	case n.Level == 2:
		styled = r.theme.H2.Render(text)
	default:
		styled = r.theme.H3.Render(text)
	}
```

With:

```go
	var styled string
	switch n.Level {
	case 1:
		styled = r.theme.H1.Render(text)
	case 2:
		styled = r.theme.H2.Render(text)
	case 3:
		styled = r.theme.H3.Render(text)
	case 4:
		styled = r.theme.H4.Render(text)
	case 5:
		styled = r.theme.H5.Render(text)
	case 6:
		styled = r.theme.H6.Render(text)
	default:
		styled = r.theme.H6.Render(text)
	}
```

- [ ] **Step 5: Run all renderer tests to verify they pass**

Run: `go test ./internal/renderer/ -v`
Expected: ALL PASS — H4/H5/H6 now render with distinct styles.

- [ ] **Step 6: Commit**

```bash
git add internal/renderer/headings.go internal/renderer/headings_test.go
git commit -m "feat: render H4, H5, H6 with distinct heading styles"
```

---

### Task 3: Update test fixture and integration tests

**Files:**
- Modify: `testdata/full.md`
- Modify: `internal/renderer/integration_test.go:43-68` (checks array)

- [ ] **Step 1: Add H4/H5/H6 headings to test fixture**

In `testdata/full.md`, add after line 9 (`### Links and Images`):

```markdown

#### Detailed Notes

Some H4-level content.

##### Implementation Details

Some H5-level content.

###### Footnotes

Some H6-level content.
```

- [ ] **Step 2: Add integration checks for H4/H5/H6**

In `internal/renderer/integration_test.go`, add three entries to the `checks` slice (after line 67, the `{"horizontal rule", "─"}` entry):

```go
		{"H4 heading", "Detailed Notes"},
		{"H5 heading", "Implementation Details"},
		{"H6 heading", "Footnotes"},
```

- [ ] **Step 3: Run integration tests**

Run: `go test ./internal/renderer/ -run 'TestFullDocumentRender' -v`
Expected: PASS — the fixture now includes H4–H6 headings and the integration test checks for them.

- [ ] **Step 4: Run full test suite**

Run: `go test ./...`
Expected: ALL PASS.

- [ ] **Step 5: Commit**

```bash
git add testdata/full.md internal/renderer/integration_test.go
git commit -m "test: add H4-H6 headings to test fixture and integration tests"
```
