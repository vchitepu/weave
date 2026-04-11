# Multi-File Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Allow `weave file1.md file2.md ...` to render multiple Markdown files sequentially, separated by a thematic break + filename header.

**Architecture:** Minimal changes to `cmd/weave/main.go`. A `renderFile()` helper reads and renders a single file via the existing goldmark pipeline. A `fileSeparator()` helper produces the styled separator string. The `run()` function gets a multi-file branch that loops over args, accumulates output, and makes a single paging decision at the end. No changes to renderer, pager, or theme packages.

**Tech Stack:** Go, cobra, goldmark, lipgloss, golang.org/x/term (all existing dependencies)

---

### Task 1: Add `fileSeparator()` helper with tests

**Files:**
- Modify: `cmd/weave/main.go` (add `fileSeparator` function after `normalizeWidth` at line 132)
- Modify: `cmd/weave/main_test.go` (add tests)

- [ ] **Step 1: Write the failing test for `fileSeparator()`**

Add to `cmd/weave/main_test.go`:

```go
func TestFileSeparator_ContainsFilename(t *testing.T) {
	th := theme.DarkTheme()
	sep := fileSeparator("notes.md", 80, th)
	if !strings.Contains(sep, "notes.md") {
		t.Fatalf("separator should contain filename, got: %q", sep)
	}
}

func TestFileSeparator_ContainsRule(t *testing.T) {
	th := theme.DarkTheme()
	sep := fileSeparator("notes.md", 80, th)
	if !strings.Contains(sep, "─") {
		t.Fatalf("separator should contain a rule character, got: %q", sep)
	}
}

func TestFileSeparator_EndsWithNewline(t *testing.T) {
	th := theme.DarkTheme()
	sep := fileSeparator("notes.md", 80, th)
	if !strings.HasSuffix(sep, "\n") {
		t.Fatalf("separator should end with newline, got: %q", sep)
	}
}
```

You will also need to add `"strings"` and `"github.com/vchitepu/weave/internal/theme"` to the import block in `main_test.go`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/weave/ -run TestFileSeparator -v`

Expected: FAIL — `fileSeparator` is undefined.

- [ ] **Step 3: Implement `fileSeparator()`**

Add to the bottom of `cmd/weave/main.go` (after `normalizeWidth`):

```go
// fileSeparator returns a visual separator placed before a file's rendered
// output in multi-file mode. It consists of a full-width thematic break
// followed by the filename on its own line, styled as a dim label.
func fileSeparator(filename string, width int, th theme.Theme) string {
	contentWidth := width - rightMargin - leftPad
	if contentWidth < 20 {
		contentWidth = 20
	}
	rule := th.HorizontalRule.Render(strings.Repeat("─", contentWidth))
	label := th.Dim.Render(filename)
	return "\n" + pad + rule + "\n" + pad + label + "\n\n"
}
```

You will also need to add the following imports to `main.go` (they are not yet present):
- `"github.com/vchitepu/weave/internal/renderer"` is already imported (uses `renderer.Priority`)
- The `pad`, `leftPad`, and `rightMargin` constants are in the `renderer` package. Since `fileSeparator` lives in `main.go` (package `main`), you need to reference them. However, `pad`, `leftPad`, and `rightMargin` are unexported in `internal/renderer`. Instead, define local constants in `main.go`:

```go
const (
	separatorLeftPad    = 2
	separatorRightMargin = 2
	separatorPad        = "  "
)
```

And use those in `fileSeparator` instead of `pad`, `leftPad`, and `rightMargin`:

```go
func fileSeparator(filename string, width int, th theme.Theme) string {
	contentWidth := width - separatorRightMargin - separatorLeftPad
	if contentWidth < 20 {
		contentWidth = 20
	}
	rule := th.HorizontalRule.Render(strings.Repeat("─", contentWidth))
	label := th.Dim.Render(filename)
	return "\n" + separatorPad + rule + "\n" + separatorPad + label + "\n\n"
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/weave/ -run TestFileSeparator -v`

Expected: All 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/weave/main.go cmd/weave/main_test.go
git commit -m "feat: add fileSeparator helper for multi-file output"
```

---

### Task 2: Add `renderFile()` helper with tests

**Files:**
- Modify: `cmd/weave/main.go` (add `renderFile` function)
- Modify: `cmd/weave/main_test.go` (add tests)

- [ ] **Step 1: Write the failing tests for `renderFile()`**

Add to `cmd/weave/main_test.go`:

```go
func TestRenderFile_ValidFile(t *testing.T) {
	// Create a temp markdown file
	tmpFile, err := os.CreateTemp("", "weave-test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("# Hello\n\nWorld\n"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	md := buildMarkdown(theme.DarkTheme(), 80)
	output, err := renderFile(tmpFile.Name(), md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(output, "Hello") {
		t.Fatalf("output should contain 'Hello', got: %q", output)
	}
}

func TestRenderFile_MissingFile(t *testing.T) {
	md := buildMarkdown(theme.DarkTheme(), 80)
	_, err := renderFile("/tmp/does-not-exist-weave-test.md", md)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
```

You will also need to add `"os"` to the imports in `main_test.go`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/weave/ -run TestRenderFile -v`

Expected: FAIL — `renderFile` and `buildMarkdown` are undefined.

- [ ] **Step 3: Implement `buildMarkdown()` and `renderFile()`**

Add to `cmd/weave/main.go` (after `fileSeparator`):

```go
// buildMarkdown constructs the goldmark Markdown instance with the weave
// renderer, theme, and width baked in. Extracted so it can be shared between
// single-file, multi-file, and test paths.
func buildMarkdown(th theme.Theme, width int) goldmark.Markdown {
	r := renderer.New(th, width)
	return goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(r, renderer.Priority),
				),
			),
		),
	)
}

// renderFile reads a Markdown file from disk and renders it using the
// provided goldmark instance. Returns the rendered ANSI string or an error.
func renderFile(path string, md goldmark.Markdown) (string, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("weave: no such file: %s", path)
	}
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return "", fmt.Errorf("weave: render error for %s: %w", path, err)
	}
	return buf.String(), nil
}
```

Also update the existing code in `run()` (lines 85–96) to use `buildMarkdown`:

Replace:
```go
	// Build goldmark with our renderer
	r := renderer.New(th, width)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(r, renderer.Priority),
				),
			),
		),
	)
```

With:
```go
	// Build goldmark with our renderer
	md := buildMarkdown(th, width)
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/weave/ -run TestRenderFile -v`

Expected: Both tests PASS.

- [ ] **Step 5: Run all existing tests to verify no regressions**

Run: `go test ./...`

Expected: All tests PASS. The `buildMarkdown` refactor produces the identical goldmark instance.

- [ ] **Step 6: Commit**

```bash
git add cmd/weave/main.go cmd/weave/main_test.go
git commit -m "feat: add renderFile and buildMarkdown helpers"
```

---

### Task 3: Wire up multi-file loop in `run()`

**Files:**
- Modify: `cmd/weave/main.go` (update `run()` function and cobra args)

- [ ] **Step 1: Change cobra args constraint**

In `cmd/weave/main.go`, change line 32:

From:
```go
		Args:    cobra.MaximumNArgs(1),
```

To:
```go
		Args:    cobra.ArbitraryArgs,
```

Also update the `Use` string on line 29:

From:
```go
		Use:     "weave [file]",
```

To:
```go
		Use:     "weave [file...]",
```

- [ ] **Step 2: Add multi-file branch in `run()`**

Replace the input-reading block in `run()` (lines 45–63, from `// Read input` through the closing `}` of the else block) with the following:

```go
	// Multi-file mode: render each file with separators
	if len(args) >= 2 {
		// Detect terminal width
		width := widthFlag
		autoWidth := widthFlag == 0
		if width == 0 {
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				width = w
			} else {
				width = 80
			}
		}
		width = normalizeWidth(width, autoWidth)

		// Validate theme flag
		if themeFlag != "" && themeFlag != "dark" && themeFlag != "light" {
			return fmt.Errorf("weave: invalid theme %q (use 'dark' or 'light')", themeFlag)
		}
		th := theme.Detect(themeFlag)
		md := buildMarkdown(th, width)

		var combined strings.Builder
		for i, path := range args {
			if i > 0 {
				combined.WriteString(fileSeparator(path, width, th))
			}
			rendered, err := renderFile(path, md)
			if err != nil {
				return err
			}
			combined.WriteString(rendered)
		}

		output := combined.String()

		// Paging decision
		isTTY := term.IsTerminal(int(os.Stdout.Fd()))
		if isTTY {
			_, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				termHeight = 24
			}
			lineCount := strings.Count(output, "\n")
			if pager.ShouldPage(lineCount, termHeight) {
				return pager.Run(output)
			}
		}
		_, err := fmt.Fprint(os.Stdout, output)
		return err
	}

	// Single-file or stdin mode (existing behavior)
	var input []byte
	var err error

	if len(args) == 1 {
		input, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("weave: no such file: %s", args[0])
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return cmd.Help()
		}
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("weave: failed to read stdin: %w", err)
		}
	}
```

The rest of `run()` (width detection, theme detection, rendering, paging) remains unchanged for the single-file/stdin path.

- [ ] **Step 3: Verify the build compiles**

Run: `go build ./cmd/weave/`

Expected: Compiles without errors.

- [ ] **Step 4: Run all existing tests to verify no regressions**

Run: `go test ./...`

Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/weave/main.go
git commit -m "feat: wire up multi-file rendering loop in run()"
```

---

### Task 4: Add multi-file integration test

**Files:**
- Modify: `cmd/weave/main_test.go` (add integration-style test)

- [ ] **Step 1: Write the multi-file output test**

Add to `cmd/weave/main_test.go`:

```go
func TestMultiFileOutput(t *testing.T) {
	// Create two temp markdown files
	tmp1, err := os.CreateTemp("", "weave-multi1-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp1.Name())
	if _, err := tmp1.WriteString("# First\n\nAlpha content\n"); err != nil {
		t.Fatal(err)
	}
	tmp1.Close()

	tmp2, err := os.CreateTemp("", "weave-multi2-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp2.Name())
	if _, err := tmp2.WriteString("# Second\n\nBeta content\n"); err != nil {
		t.Fatal(err)
	}
	tmp2.Close()

	th := theme.DarkTheme()
	md := buildMarkdown(th, 80)

	// Render both files with separator
	var combined strings.Builder
	paths := []string{tmp1.Name(), tmp2.Name()}
	for i, path := range paths {
		if i > 0 {
			combined.WriteString(fileSeparator(path, 80, th))
		}
		rendered, err := renderFile(path, md)
		if err != nil {
			t.Fatalf("renderFile(%s): %v", path, err)
		}
		combined.WriteString(rendered)
	}

	output := combined.String()

	if !strings.Contains(output, "First") {
		t.Error("output should contain 'First'")
	}
	if !strings.Contains(output, "Alpha") {
		t.Error("output should contain 'Alpha'")
	}
	if !strings.Contains(output, "Second") {
		t.Error("output should contain 'Second'")
	}
	if !strings.Contains(output, "Beta") {
		t.Error("output should contain 'Beta'")
	}
	// Separator should contain the second filename and a rule
	if !strings.Contains(output, "─") {
		t.Error("output should contain a rule character in the separator")
	}
	// The separator should contain the second file's name
	base2 := tmp2.Name()
	if !strings.Contains(output, base2) {
		t.Errorf("output should contain second filename %q in separator", base2)
	}
}
```

- [ ] **Step 2: Run the test to verify it passes**

Run: `go test ./cmd/weave/ -run TestMultiFileOutput -v`

Expected: PASS.

- [ ] **Step 3: Run all tests**

Run: `go test ./...`

Expected: All tests PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/weave/main_test.go
git commit -m "test: add multi-file rendering integration test"
```

---

### Task 5: Update documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update usage section**

In `README.md`, after the existing usage examples (the code block ending with `weave --width 100 README.md`), the `weave --width 100 README.md` line should be followed by a new example. Add these lines inside the existing usage code block, before the closing triple-backtick:

```sh
# Render multiple files with separators
weave file1.md file2.md file3.md
```

- [ ] **Step 2: Update the flags table**

No new flags are needed. No changes to the flags table.

- [ ] **Step 3: Move multi-file from Future Directions to Features**

In the Features list (the bullet list under `## Features`), add:

```markdown
- Multiple-file rendering with section separators between files
```

In the `### UX Improvements` section under `## Future Directions`, remove the line:

```markdown
- Multiple-file rendering (`weave file1.md file2.md`) with section separators
```

- [ ] **Step 4: Update the `Use` string note**

No code change needed — this was already done in Task 3, Step 1 (`weave [file...]`).

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "docs: update README for multi-file support"
```

---

### Task 6: Final verification

- [ ] **Step 1: Run the full test suite**

Run: `go test ./...`

Expected: All tests PASS.

- [ ] **Step 2: Build the binary**

Run: `go build -o /tmp/weave-test ./cmd/weave`

Expected: Compiles without errors.

- [ ] **Step 3: Manual smoke test with two files**

Run: `/tmp/weave-test README.md testdata/full.md`

Expected: Both files render in sequence. A thematic break + filename label appears between them. Output pages if it exceeds terminal height.

- [ ] **Step 4: Manual smoke test — single file still works**

Run: `/tmp/weave-test README.md`

Expected: Identical behavior to before the changes.

- [ ] **Step 5: Manual smoke test — stdin still works**

Run: `echo "# Test" | /tmp/weave-test`

Expected: Renders the heading as before.

- [ ] **Step 6: Clean up temp binary**

Run: `rm /tmp/weave-test`
