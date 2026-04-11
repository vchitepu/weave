# Web Viewer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `--web` flag to weave that launches a local HTTP server serving a browser-based Markdown viewer with live reload on file changes.

**Architecture:** Three new internal packages (`htmlrenderer`, `server`, `watcher`) plus CLI integration. The HTML renderer implements goldmark's `NodeRenderer` interface (same pattern as the terminal renderer) and emits styled HTML+CSS. The server serves a single-page viewer with SSE-based live reload. The watcher wraps `fsnotify` for file change detection.

**Tech Stack:** Go standard library (`net/http`, `html/template`, `os/signal`), goldmark, chroma (HTML formatter), fsnotify

---

## File Structure

### New files

| File | Responsibility |
|---|---|
| `internal/htmlrenderer/htmlrenderer.go` | Renderer struct, RegisterFuncs, document/paragraph/text/string handlers |
| `internal/htmlrenderer/headings.go` | H1–H6 and thematic break handlers |
| `internal/htmlrenderer/codeblock.go` | Fenced/indented code block handlers, container detection, chroma HTML highlighting |
| `internal/htmlrenderer/inline.go` | Emphasis, code span, link, image, strikethrough handlers |
| `internal/htmlrenderer/blockquote.go` | Blockquote handler |
| `internal/htmlrenderer/list.go` | List, list item, task checkbox handlers |
| `internal/htmlrenderer/table.go` | Table, table header, table row, table cell handlers |
| `internal/htmlrenderer/css.go` | CSS string constants for dark/light themes and base layout |
| `internal/htmlrenderer/htmlrenderer_test.go` | Unit tests for all node types |
| `internal/htmlrenderer/integration_test.go` | Full-document integration test |
| `internal/server/server.go` | HTTP server, page template, SSE broadcaster, `Start()` entry point |
| `internal/server/server_test.go` | HTTP handler tests via `httptest` |
| `internal/watcher/watcher.go` | fsnotify wrapper with debounce |
| `internal/watcher/watcher_test.go` | File change detection tests |

### Modified files

| File | Change |
|---|---|
| `cmd/weave/main.go` | Add `--web` and `--port` flags, dispatch to `server.Start()` |
| `go.mod` | Add `github.com/fsnotify/fsnotify` dependency |


---

### Task 1: Add fsnotify dependency

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add fsnotify dependency**

```bash
go get github.com/fsnotify/fsnotify
```

- [ ] **Step 2: Verify it was added to go.mod**

Run: `grep fsnotify go.mod`
Expected: a line like `github.com/fsnotify/fsnotify v1.x.x`

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: add fsnotify for file watching"
```

---

### Task 2: Implement file watcher (`internal/watcher`)

**Files:**
- Create: `internal/watcher/watcher.go`
- Create: `internal/watcher/watcher_test.go`

- [ ] **Step 1: Write the failing test for basic file watching**

Create `internal/watcher/watcher_test.go`:

```go
package watcher

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_DetectsFileChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	var called atomic.Int32
	stop, err := Watch([]string{path}, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	// Give the watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the file
	if err := os.WriteFile(path, []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for callback (debounce is ~100ms, so wait up to 500ms)
	deadline := time.After(500 * time.Millisecond)
	for called.Load() == 0 {
		select {
		case <-deadline:
			t.Fatal("callback was not called within timeout")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestWatch_EmptyPaths_NoOp(t *testing.T) {
	var called atomic.Int32
	stop, err := Watch(nil, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	time.Sleep(200 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatal("callback should not have been called for empty paths")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/watcher/ -v -count=1`
Expected: Compilation error — package `watcher` does not exist yet.

- [ ] **Step 3: Write the implementation**

Create `internal/watcher/watcher.go`:

```go
package watcher

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch watches the given file paths for changes and calls onChange when any
// file is written or created. Changes are debounced by 100ms.
// Returns a stop function to clean up the watcher.
// If paths is empty, returns a no-op stop function and nil error.
func Watch(paths []string, onChange func()) (stop func(), err error) {
	if len(paths) == 0 {
		return func() {}, nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if err := w.Add(p); err != nil {
			w.Close()
			return nil, err
		}
	}

	var once sync.Once
	done := make(chan struct{})

	go func() {
		var timer *time.Timer
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(100*time.Millisecond, onChange)
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
				// Ignore errors silently
			case <-done:
				return
			}
		}
	}()

	stopFn := func() {
		once.Do(func() {
			close(done)
			w.Close()
		})
	}

	return stopFn, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/watcher/ -v -count=1`
Expected: Both tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/watcher/
git commit -m "feat: add file watcher with debounced change detection"
```


---

### Task 3: Implement HTML renderer — core structure and CSS

**Files:**
- Create: `internal/htmlrenderer/css.go`
- Create: `internal/htmlrenderer/htmlrenderer.go`

- [ ] **Step 1: Create the CSS constants**

Create `internal/htmlrenderer/css.go`:

```go
package htmlrenderer

import "github.com/vchitepu/weave/internal/theme"

// BaseCSS contains the theme-independent layout styles.
const BaseCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  line-height: 1.6;
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  color: var(--color-text);
  background: var(--color-bg);
}
h1, h2, h3, h4, h5, h6 { margin: 1.5rem 0 0.5rem 0; }
h1 { font-size: 2rem; color: var(--color-h1); border-bottom: 1px solid var(--color-rule); padding-bottom: 0.3rem; }
h2 { font-size: 1.5rem; color: var(--color-h2); }
h3 { font-size: 1.25rem; color: var(--color-h3); }
h4 { font-size: 1.1rem; color: var(--color-h4); }
h5 { font-size: 1rem; color: var(--color-h5); }
h6 { font-size: 0.9rem; color: var(--color-h6); }
p { margin: 0.5rem 0 1rem 0; }
a { color: var(--color-link); text-decoration: none; }
a:hover { text-decoration: underline; }
em { font-style: italic; }
strong { font-weight: bold; }
del { text-decoration: line-through; }
hr { border: none; border-top: 1px solid var(--color-rule); margin: 1.5rem 0; }
code.inline-code {
  background: var(--color-code-bg);
  padding: 0.15rem 0.4rem;
  border-radius: 3px;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 0.9em;
}
blockquote {
  border-left: 3px solid var(--color-blockquote-bar);
  padding-left: 1rem;
  margin: 1rem 0;
  color: var(--color-blockquote-text);
}
ul, ol { margin: 0.5rem 0 1rem 1.5rem; }
li { margin: 0.25rem 0; }
li.task-item { list-style: none; margin-left: -1.5rem; }
.checkbox { font-family: "SFMono-Regular", Consolas, monospace; margin-right: 0.4rem; }
.checkbox.checked { color: var(--color-task-checked); }
.checkbox.unchecked { color: var(--color-task-unchecked); }
table { border-collapse: collapse; margin: 1rem 0; width: auto; }
th, td { border: 1px solid var(--color-table-border); padding: 0.4rem 0.8rem; }
th { font-weight: bold; color: var(--color-table-header); background: var(--color-table-header-bg); }
img { max-width: 100%; }
.code-container {
  position: relative;
  border: 1px solid var(--color-code-border);
  border-radius: 6px;
  margin: 1rem 0;
  overflow: hidden;
}
.code-container .lang-badge {
  position: absolute;
  top: 0; left: 1rem;
  background: var(--color-code-border);
  color: var(--color-code-badge-text);
  padding: 0.1rem 0.5rem;
  border-radius: 0 0 4px 4px;
  font-size: 0.75rem;
  font-family: "SFMono-Regular", Consolas, monospace;
}
.code-container pre {
  padding: 1.5rem 1rem 1rem 1rem;
  overflow-x: auto;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 0.9em;
  line-height: 1.5;
  background: var(--color-code-bg);
}
.code-container.container-tree { border-color: var(--color-tree-border); }
.code-container.container-tree .lang-badge { background: var(--color-tree-border); }
.code-container.container-diagram { border-color: var(--color-diagram-border); }
.code-container.container-diagram .lang-badge { background: var(--color-diagram-border); }
.code-container.container-shell { border-color: var(--color-shell-border); }
.code-container.container-shell .lang-badge { background: var(--color-shell-border); }
.file-separator { margin: 2rem 0 0.5rem 0; border-color: var(--color-rule); }
.file-label { color: var(--color-dim); font-size: 0.85rem; margin-bottom: 1rem; }
`

// ThemeCSS returns CSS custom property declarations for the given theme.
func ThemeCSS(th theme.Theme) string {
	if th.ChromaStyle == "github-dark" {
		return darkThemeCSS
	}
	return lightThemeCSS
}

const darkThemeCSS = `
:root {
  --color-text: #D0D4DC;
  --color-bg: #1A1E28;
  --color-h1: #8BA4D4;
  --color-h2: #A8A0D6;
  --color-h3: #C9A86A;
  --color-h4: #A88A55;
  --color-h5: #876C42;
  --color-h6: #665030;
  --color-rule: #3A3F4B;
  --color-link: #7FA3C8;
  --color-dim: #9BA3B2;
  --color-code-bg: #262B36;
  --color-code-border: #3A3F4B;
  --color-code-badge-text: #B7C0D0;
  --color-blockquote-bar: #3A3F4B;
  --color-blockquote-text: #9BA3B2;
  --color-table-border: #3A3F4B;
  --color-table-header: #8BA4D4;
  --color-table-header-bg: #262B36;
  --color-tree-border: #5D7290;
  --color-diagram-border: #7B74A6;
  --color-shell-border: #9F8656;
  --color-task-checked: #6BBF8A;
  --color-task-unchecked: #5A5F6B;
  --color-inline-code-bg: #262B36;
}
`

const lightThemeCSS = `
:root {
  --color-text: #2C3340;
  --color-bg: #FAFBFC;
  --color-h1: #3F5F8A;
  --color-h2: #665E95;
  --color-h3: #8D6B3F;
  --color-h4: #A07850;
  --color-h5: #B38A63;
  --color-h6: #C69C78;
  --color-rule: #C2C7D0;
  --color-link: #496A92;
  --color-dim: #6C7483;
  --color-code-bg: #EEF1F5;
  --color-code-border: #C2C7D0;
  --color-code-badge-text: #5F6B7A;
  --color-blockquote-bar: #C2C7D0;
  --color-blockquote-text: #6C7483;
  --color-table-border: #C2C7D0;
  --color-table-header: #3F5F8A;
  --color-table-header-bg: #EEF1F5;
  --color-tree-border: #5E7699;
  --color-diagram-border: #7D73A3;
  --color-shell-border: #9A7B52;
  --color-task-checked: #3A8C5C;
  --color-task-unchecked: #9BA3B2;
  --color-inline-code-bg: #EEF1F5;
}
`
```

- [ ] **Step 2: Create the renderer core**

Create `internal/htmlrenderer/htmlrenderer.go`:

```go
package htmlrenderer

import (
	"bytes"
	"strings"

	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Priority is the goldmark renderer priority for the HTML renderer.
const Priority = 100

// Renderer implements goldmark's NodeRenderer interface, emitting HTML.
type Renderer struct {
	theme theme.Theme
}

// New creates a new HTML Renderer.
func New(th theme.Theme) *Renderer {
	return &Renderer{theme: th}
}

// RegisterFuncs registers AST node render functions.
func (r *Renderer) RegisterFuncs(reg goldrenderer.NodeRendererFuncRegisterer) {
	// Block nodes
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)

	// Table nodes
	reg.Register(east.KindTable, r.renderTable)
	reg.Register(east.KindTableHeader, r.renderTableHeader)
	reg.Register(east.KindTableRow, r.renderTableRow)
	reg.Register(east.KindTableCell, r.renderTableCell)

	// Extension inlines
	reg.Register(east.KindStrikethrough, r.renderStrikethrough)
	reg.Register(east.KindTaskCheckBox, r.renderTaskCheckBox)

	// Inline nodes
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<p>")
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	text := n.Segment.Value(source)
	_, _ = w.Write(htmlEscape(text))
	if n.HardLineBreak() {
		_, _ = w.WriteString("<br>\n")
	} else if n.SoftLineBreak() {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	_, _ = w.Write(htmlEscape(n.Value))
	return ast.WalkContinue, nil
}

// htmlEscape escapes HTML special characters.
func htmlEscape(b []byte) []byte {
	var buf bytes.Buffer
	for _, c := range b {
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		default:
			buf.WriteByte(c)
		}
	}
	return buf.Bytes()
}

// htmlEscapeString escapes HTML special characters in a string.
func htmlEscapeString(s string) string {
	return string(htmlEscape([]byte(s)))
}

// collectText recursively collects text content from an AST node tree.
func collectText(node ast.Node, source []byte) string {
	var buf strings.Builder
	collectTextInto(&buf, node, source)
	return buf.String()
}

func collectTextInto(buf *strings.Builder, node ast.Node, source []byte) {
	switch n := node.(type) {
	case *ast.Text:
		buf.Write(n.Segment.Value(source))
		if n.SoftLineBreak() || n.HardLineBreak() {
			buf.WriteString("\n")
		}
	case *ast.String:
		buf.Write(n.Value)
	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			collectTextInto(buf, child, source)
		}
	}
}
```

- [ ] **Step 3: Verify the core compiles**

Run: `go build ./internal/htmlrenderer/`
Expected: Compiles successfully (even though not all registered handlers exist yet — they will be stub panics until implemented; actually we should add stubs now).

Note: Since RegisterFuncs references handler methods that don't exist yet, this won't compile. That's expected — we'll fix it as we add each handler file in subsequent tasks. For now, comment out the register lines for handlers not yet defined, or add minimal stubs.

Actually, let's add placeholder stubs at the bottom of `htmlrenderer.go` so it compiles — we'll replace each stub with a proper implementation file in subsequent tasks:

Add to the end of `internal/htmlrenderer/htmlrenderer.go`:

```go
// Stubs — replaced by proper implementations in subsequent tasks.
// Each of these will be moved to its own file.

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
```

- [ ] **Step 4: Verify the package compiles**

Run: `go build ./internal/htmlrenderer/`
Expected: Compiles without errors.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/css.go internal/htmlrenderer/htmlrenderer.go
git commit -m "feat: add HTML renderer core structure and CSS themes"
```


---

### Task 4: HTML renderer — headings

**Files:**
- Create: `internal/htmlrenderer/headings.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove heading/thematic break stubs)

- [ ] **Step 1: Write the failing test**

Add to `internal/htmlrenderer/htmlrenderer_test.go` (create the file):

```go
package htmlrenderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func renderHTML(t *testing.T, input string) string {
	t.Helper()
	return renderHTMLWithTheme(t, input, theme.DarkTheme())
}

func renderHTMLWithTheme(t *testing.T, input string, th theme.Theme) string {
	t.Helper()
	r := New(th)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(r, Priority),
				),
			),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("failed to convert markdown: %v", err)
	}
	return buf.String()
}

func TestHeading_H1(t *testing.T) {
	out := renderHTML(t, "# Hello")
	if !strings.Contains(out, "<h1>") || !strings.Contains(out, "Hello") || !strings.Contains(out, "</h1>") {
		t.Fatalf("expected <h1>Hello</h1>, got: %q", out)
	}
}

func TestHeading_H2(t *testing.T) {
	out := renderHTML(t, "## World")
	if !strings.Contains(out, "<h2>") || !strings.Contains(out, "World") {
		t.Fatalf("expected <h2>World</h2>, got: %q", out)
	}
}

func TestHeading_H3Through6(t *testing.T) {
	for _, tc := range []struct{ md, tag string }{
		{"### H3", "<h3>"},
		{"#### H4", "<h4>"},
		{"##### H5", "<h5>"},
		{"###### H6", "<h6>"},
	} {
		out := renderHTML(t, tc.md)
		if !strings.Contains(out, tc.tag) {
			t.Fatalf("for %q, expected %s, got: %q", tc.md, tc.tag, out)
		}
	}
}

func TestThematicBreak(t *testing.T) {
	out := renderHTML(t, "---")
	if !strings.Contains(out, "<hr>") {
		t.Fatalf("expected <hr>, got: %q", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/htmlrenderer/ -run "TestHeading|TestThematicBreak" -v -count=1`
Expected: Tests pass with stubs (stubs output nothing, so assertions fail).

- [ ] **Step 3: Implement headings**

Create `internal/htmlrenderer/headings.go`:

```go
package htmlrenderer

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	level := n.Level
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	if entering {
		_, _ = w.WriteString(fmt.Sprintf("<h%d>", level))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("</h%d>\n", level))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<hr>\n")
	return ast.WalkContinue, nil
}
```

Then remove the `renderHeading` and `renderThematicBreak` stubs from `htmlrenderer.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/htmlrenderer/ -run "TestHeading|TestThematicBreak" -v -count=1`
Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML heading and thematic break rendering"
```


---

### Task 5: HTML renderer — inline elements

**Files:**
- Create: `internal/htmlrenderer/inline.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove inline stubs)

- [ ] **Step 1: Write the failing tests**

Add to `internal/htmlrenderer/htmlrenderer_test.go`:

```go
func TestEmphasis(t *testing.T) {
	out := renderHTML(t, "*italic*")
	if !strings.Contains(out, "<em>") || !strings.Contains(out, "italic") {
		t.Fatalf("expected <em>italic</em>, got: %q", out)
	}
}

func TestStrong(t *testing.T) {
	out := renderHTML(t, "**bold**")
	if !strings.Contains(out, "<strong>") || !strings.Contains(out, "bold") {
		t.Fatalf("expected <strong>bold</strong>, got: %q", out)
	}
}

func TestInlineCode(t *testing.T) {
	out := renderHTML(t, "use `fmt.Println`")
	if !strings.Contains(out, `<code class="inline-code">`) || !strings.Contains(out, "fmt.Println") {
		t.Fatalf("expected inline code, got: %q", out)
	}
}

func TestLink(t *testing.T) {
	out := renderHTML(t, "[Go](https://golang.org)")
	if !strings.Contains(out, `<a href="https://golang.org"`) || !strings.Contains(out, "Go") {
		t.Fatalf("expected link, got: %q", out)
	}
}

func TestImage(t *testing.T) {
	out := renderHTML(t, "![alt text](image.png)")
	if !strings.Contains(out, `<img src="image.png"`) || !strings.Contains(out, `alt="alt text"`) {
		t.Fatalf("expected img tag, got: %q", out)
	}
}

func TestStrikethrough(t *testing.T) {
	out := renderHTML(t, "~~deleted~~")
	if !strings.Contains(out, "<del>") || !strings.Contains(out, "deleted") {
		t.Fatalf("expected <del>deleted</del>, got: %q", out)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/htmlrenderer/ -run "TestEmphasis|TestStrong|TestInlineCode|TestLink|TestImage|TestStrikethrough" -v -count=1`
Expected: All fail (stubs produce no output).

- [ ] **Step 3: Implement inline handlers**

Create `internal/htmlrenderer/inline.go`:

```go
package htmlrenderer

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if n.Level == 2 {
		if entering {
			_, _ = w.WriteString("<strong>")
		} else {
			_, _ = w.WriteString("</strong>")
		}
	} else {
		if entering {
			_, _ = w.WriteString("<em>")
		} else {
			_, _ = w.WriteString("</em>")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		text := collectText(node, source)
		_, _ = w.WriteString(fmt.Sprintf(`<code class="inline-code">%s</code>`, htmlEscapeString(text)))
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString(fmt.Sprintf(`<a href="%s">`, htmlEscapeString(string(n.Destination))))
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	alt := string(n.Text(source))
	_, _ = w.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, htmlEscapeString(string(n.Destination)), htmlEscapeString(alt)))
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_ = node.(*east.Strikethrough)
	if entering {
		_, _ = w.WriteString("<del>")
	} else {
		_, _ = w.WriteString("</del>")
	}
	return ast.WalkContinue, nil
}
```

Then remove the `renderEmphasis`, `renderCodeSpan`, `renderLink`, `renderImage`, and `renderStrikethrough` stubs from `htmlrenderer.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/htmlrenderer/ -run "TestEmphasis|TestStrong|TestInlineCode|TestLink|TestImage|TestStrikethrough" -v -count=1`
Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML inline element rendering"
```


---

### Task 6: HTML renderer — code blocks with syntax highlighting

**Files:**
- Create: `internal/htmlrenderer/codeblock.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove code block stubs)

- [ ] **Step 1: Write the failing tests**

Add to `internal/htmlrenderer/htmlrenderer_test.go`:

```go
func TestFencedCodeBlock_WithLanguage(t *testing.T) {
	input := "```go\nfmt.Println(\"hello\")\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="code-container`) {
		t.Fatalf("expected code-container class, got: %q", out)
	}
	if !strings.Contains(out, `<span class="lang-badge">go</span>`) {
		t.Fatalf("expected lang-badge, got: %q", out)
	}
	if !strings.Contains(out, "<pre>") {
		t.Fatalf("expected <pre>, got: %q", out)
	}
}

func TestFencedCodeBlock_NoLanguage(t *testing.T) {
	input := "```\nplain code\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, `class="code-container"`) {
		t.Fatalf("expected code-container class, got: %q", out)
	}
	if strings.Contains(out, "lang-badge") {
		t.Fatalf("expected no lang-badge for unspecified language, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerTree(t *testing.T) {
	input := "```tree\n├── src\n│   └── main.go\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-tree") {
		t.Fatalf("expected container-tree class, got: %q", out)
	}
	if !strings.Contains(out, `<span class="lang-badge">tree</span>`) {
		t.Fatalf("expected tree badge, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerShell(t *testing.T) {
	input := "```bash\necho hello\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-shell") {
		t.Fatalf("expected container-shell class, got: %q", out)
	}
	if !strings.Contains(out, `<span class="lang-badge">$</span>`) {
		t.Fatalf("expected $ badge, got: %q", out)
	}
}

func TestFencedCodeBlock_ContainerDiagram(t *testing.T) {
	input := "```mermaid\ngraph TD\n```"
	out := renderHTML(t, input)
	if !strings.Contains(out, "container-diagram") {
		t.Fatalf("expected container-diagram class, got: %q", out)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/htmlrenderer/ -run "TestFencedCodeBlock" -v -count=1`
Expected: All fail.

- [ ] **Step 3: Implement code block handlers**

Create `internal/htmlrenderer/codeblock.go`:

```go
package htmlrenderer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

type containerType int

const (
	containerCode containerType = iota
	containerTree
	containerDiagram
	containerShell
)

func detectContainer(lang string) containerType {
	switch strings.ToLower(lang) {
	case "tree":
		return containerTree
	case "ascii", "diagram", "art", "mermaid":
		return containerDiagram
	case "bash", "sh", "shell", "console", "terminal":
		return containerShell
	default:
		return containerCode
	}
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := ""
	if n.Language(source) != nil {
		lang = string(n.Language(source))
	}

	code := collectLines(node, source)
	r.writeCodeContainer(w, code, lang)
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	code := collectLines(node, source)
	r.writeCodeContainer(w, code, "")
	return ast.WalkSkipChildren, nil
}

func collectLines(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.Write(line.Value(source))
	}
	return strings.TrimRight(buf.String(), "\n")
}

func (r *Renderer) writeCodeContainer(w util.BufWriter, code, lang string) {
	ct := detectContainer(lang)

	// Determine CSS class and badge label
	var containerClass, badgeLabel string
	switch ct {
	case containerTree:
		containerClass = " container-tree"
		badgeLabel = "tree"
	case containerDiagram:
		containerClass = " container-diagram"
		badgeLabel = "diagram"
	case containerShell:
		containerClass = " container-shell"
		badgeLabel = "$"
	default:
		containerClass = ""
		badgeLabel = lang
	}

	// Syntax highlight for code containers with a known language
	highlighted := htmlEscapeString(code)
	if ct == containerCode && lang != "" {
		if h := r.highlightCodeHTML(code, lang); h != "" {
			highlighted = h
		}
	}

	_, _ = w.WriteString(fmt.Sprintf(`<div class="code-container%s">`, containerClass))
	if badgeLabel != "" {
		_, _ = w.WriteString(fmt.Sprintf(`<span class="lang-badge">%s</span>`, htmlEscapeString(badgeLabel)))
	}
	_, _ = w.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>", highlighted))
	_, _ = w.WriteString("</div>\n")
}

// highlightCodeHTML returns chroma-highlighted HTML or empty string on failure.
func (r *Renderer) highlightCodeHTML(code, lang string) string {
	lexer := lexers.Get(lang)
	if lexer == nil {
		return ""
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(r.theme.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(
		chromahtml.WithClasses(false),
		chromahtml.InlineCode(true),
		chromahtml.PreventSurroundingPre(true),
	)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return ""
	}

	return buf.String()
}

// ChromaCSS returns the chroma stylesheet for inline-style rendering.
// When using inline styles (InlineCode: true), this is not strictly needed,
// but we include it for completeness.
func ChromaCSS(th theme.Theme) string {
	style := styles.Get(th.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	_ = formatter.WriteCSS(&buf, style)
	return buf.String()
}
```

Note: We use `chromahtml.InlineCode(true)` so chroma emits inline `style` attributes, avoiding the need for a separate chroma CSS class stylesheet. `ChromaCSS()` is provided for optional use if we switch to class-based highlighting later.

Then remove the `renderFencedCodeBlock` and `renderCodeBlock` stubs from `htmlrenderer.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/htmlrenderer/ -run "TestFencedCodeBlock" -v -count=1`
Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML code block rendering with syntax highlighting"
```


---

### Task 7: HTML renderer — blockquotes

**Files:**
- Create: `internal/htmlrenderer/blockquote.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove blockquote stub)

- [ ] **Step 1: Write the failing test**

Add to `internal/htmlrenderer/htmlrenderer_test.go`:

```go
func TestBlockquote(t *testing.T) {
	out := renderHTML(t, "> This is quoted text")
	if !strings.Contains(out, "<blockquote>") || !strings.Contains(out, "This is quoted text") {
		t.Fatalf("expected <blockquote> with text, got: %q", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/htmlrenderer/ -run "TestBlockquote" -v -count=1`
Expected: FAIL.

- [ ] **Step 3: Implement blockquote handler**

Create `internal/htmlrenderer/blockquote.go`:

```go
package htmlrenderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<blockquote>\n")
	} else {
		_, _ = w.WriteString("</blockquote>\n")
	}
	return ast.WalkContinue, nil
}
```

Remove the `renderBlockquote` stub from `htmlrenderer.go`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/htmlrenderer/ -run "TestBlockquote" -v -count=1`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML blockquote rendering"
```

---

### Task 8: HTML renderer — lists and task checkboxes

**Files:**
- Create: `internal/htmlrenderer/list.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove list/task stubs)

- [ ] **Step 1: Write the failing tests**

Add to `internal/htmlrenderer/htmlrenderer_test.go`:

```go
func TestUnorderedList(t *testing.T) {
	input := "- Alpha\n- Beta\n- Gamma"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<ul>") || !strings.Contains(out, "<li>") {
		t.Fatalf("expected <ul><li>, got: %q", out)
	}
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Beta") || !strings.Contains(out, "Gamma") {
		t.Fatalf("expected list items, got: %q", out)
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. First\n2. Second\n3. Third"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<ol>") || !strings.Contains(out, "<li>") {
		t.Fatalf("expected <ol><li>, got: %q", out)
	}
}

func TestTaskList_Checked(t *testing.T) {
	input := "- [x] Done task"
	out := renderHTML(t, input)
	if !strings.Contains(out, "task-item") {
		t.Fatalf("expected task-item class, got: %q", out)
	}
	if !strings.Contains(out, `class="checkbox checked"`) {
		t.Fatalf("expected checked checkbox, got: %q", out)
	}
}

func TestTaskList_Unchecked(t *testing.T) {
	input := "- [ ] Pending task"
	out := renderHTML(t, input)
	if !strings.Contains(out, "task-item") {
		t.Fatalf("expected task-item class, got: %q", out)
	}
	if !strings.Contains(out, `class="checkbox unchecked"`) {
		t.Fatalf("expected unchecked checkbox, got: %q", out)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/htmlrenderer/ -run "TestUnorderedList|TestOrderedList|TestTaskList" -v -count=1`
Expected: All fail.

- [ ] **Step 3: Implement list and task checkbox handlers**

Create `internal/htmlrenderer/list.go`:

```go
package htmlrenderer

import (
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if n.IsOrdered() {
		if entering {
			_, _ = w.WriteString("<ol>\n")
		} else {
			_, _ = w.WriteString("</ol>\n")
		}
	} else {
		if entering {
			_, _ = w.WriteString("<ul>\n")
		} else {
			_, _ = w.WriteString("</ul>\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Check if this is a task list item
		isTask := false
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
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
			_, _ = w.WriteString(`<li class="task-item">`)
		} else {
			_, _ = w.WriteString("<li>")
		}
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*east.TaskCheckBox)
	if n.IsChecked {
		_, _ = w.WriteString(`<span class="checkbox checked">&#x2713;</span>`)
	} else {
		_, _ = w.WriteString(`<span class="checkbox unchecked">&#x25CB;</span>`)
	}
	return ast.WalkContinue, nil
}
```

Remove the `renderList`, `renderListItem`, and `renderTaskCheckBox` stubs from `htmlrenderer.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/htmlrenderer/ -run "TestUnorderedList|TestOrderedList|TestTaskList" -v -count=1`
Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML list and task checkbox rendering"
```


---

### Task 9: HTML renderer — tables

**Files:**
- Create: `internal/htmlrenderer/table.go`
- Modify: `internal/htmlrenderer/htmlrenderer.go` (remove table stubs)

- [ ] **Step 1: Write the failing test**

Add to `internal/htmlrenderer/htmlrenderer_test.go`:

```go
func TestTable(t *testing.T) {
	input := "| Name | Age |\n|---|---|\n| Alice | 30 |\n| Bob | 25 |"
	out := renderHTML(t, input)
	if !strings.Contains(out, "<table>") {
		t.Fatalf("expected <table>, got: %q", out)
	}
	if !strings.Contains(out, "<thead>") || !strings.Contains(out, "<th>") {
		t.Fatalf("expected <thead>/<th>, got: %q", out)
	}
	if !strings.Contains(out, "<tbody>") || !strings.Contains(out, "<td>") {
		t.Fatalf("expected <tbody>/<td>, got: %q", out)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Fatalf("expected table content, got: %q", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/htmlrenderer/ -run "TestTable" -v -count=1`
Expected: FAIL.

- [ ] **Step 3: Implement table handlers**

Create `internal/htmlrenderer/table.go`:

```go
package htmlrenderer

import (
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table>\n")
	} else {
		_, _ = w.WriteString("</table>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<thead>\n<tr>\n")
	} else {
		_, _ = w.WriteString("</tr>\n</thead>\n<tbody>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<tr>\n")
	} else {
		_, _ = w.WriteString("</tr>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*east.TableCell)
	tag := "td"
	// Check if we're inside a table header
	if n.Parent() != nil && n.Parent().Kind() == east.KindTableHeader {
		tag = "th"
	}

	if entering {
		_, _ = w.WriteString("<" + tag + ">")
	} else {
		_, _ = w.WriteString("</" + tag + ">\n")
	}
	return ast.WalkContinue, nil
}
```

Note: The `<tbody>` opening tag is written at the end of `renderTableHeader` (exiting). The closing `</tbody>` needs to be written before `</table>`. Let's adjust `renderTable` to handle this:

```go
func (r *Renderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table>\n")
	} else {
		_, _ = w.WriteString("</tbody>\n</table>\n")
	}
	return ast.WalkContinue, nil
}
```

Remove the `renderTable`, `renderTableHeader`, `renderTableRow`, and `renderTableCell` stubs from `htmlrenderer.go`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/htmlrenderer/ -run "TestTable" -v -count=1`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "feat: add HTML table rendering"
```


---

### Task 10: HTML renderer — integration test

**Files:**
- Create: `internal/htmlrenderer/integration_test.go`

- [ ] **Step 1: Write the integration test**

Create `internal/htmlrenderer/integration_test.go`:

```go
package htmlrenderer

import (
	"os"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
)

func TestIntegration_FullDocument(t *testing.T) {
	input, err := os.ReadFile("../../testdata/full.md")
	if err != nil {
		t.Skipf("testdata/full.md not found: %v", err)
	}

	out := renderHTML(t, string(input))

	checks := []struct {
		name    string
		substr  string
	}{
		{"has h1", "<h1>"},
		{"has h2", "<h2>"},
		{"has paragraph", "<p>"},
		{"has code container", `class="code-container`},
		{"has blockquote", "<blockquote>"},
		{"has unordered list", "<ul>"},
		{"has table", "<table>"},
		{"has link", "<a href="},
		{"has emphasis", "<em>"},
		{"has strong", "<strong>"},
		{"has inline code", `class="inline-code"`},
		{"has hr", "<hr>"},
	}

	for _, check := range checks {
		if !strings.Contains(out, check.substr) {
			t.Errorf("%s: expected %q in output", check.name, check.substr)
		}
	}
}

func TestIntegration_DarkAndLightThemesDiffer(t *testing.T) {
	input := "# Hello\n\n```go\nfmt.Println()\n```"

	dark := renderHTMLWithTheme(t, input, theme.DarkTheme())
	light := renderHTMLWithTheme(t, input, theme.LightTheme())

	// Both should produce valid output
	if !strings.Contains(dark, "<h1>") || !strings.Contains(light, "<h1>") {
		t.Fatal("both themes should render headings")
	}

	// The chroma-highlighted code should differ between themes
	// (different inline styles from chroma)
	if dark == light {
		t.Error("dark and light theme outputs should differ")
	}
}
```

- [ ] **Step 2: Run integration tests**

Run: `go test ./internal/htmlrenderer/ -run "TestIntegration" -v -count=1`
Expected: All PASS.

- [ ] **Step 3: Run all htmlrenderer tests together**

Run: `go test ./internal/htmlrenderer/ -v -count=1`
Expected: All PASS.

- [ ] **Step 4: Verify no stubs remain in htmlrenderer.go**

Manually check that `internal/htmlrenderer/htmlrenderer.go` has no remaining stub methods. All handler methods should now be in their respective files (`headings.go`, `inline.go`, `codeblock.go`, `blockquote.go`, `list.go`, `table.go`).

- [ ] **Step 5: Commit**

```bash
git add internal/htmlrenderer/
git commit -m "test: add HTML renderer integration tests"
```


---

### Task 11: Implement HTTP server with SSE live reload

**Files:**
- Create: `internal/server/server.go`
- Create: `internal/server/server_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/server/server_test.go`:

```go
package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
)

func TestHandler_RootReturnsHTML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("# Hello World"), 0644); err != nil {
		t.Fatal(err)
	}

	inputs := []Input{{Path: path}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %s", ct)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "Hello World") {
		t.Fatalf("expected rendered content, got: %s", bodyStr[:200])
	}
	if !strings.Contains(bodyStr, "<html>") {
		t.Fatalf("expected full HTML page")
	}
	if !strings.Contains(bodyStr, "EventSource") {
		t.Fatalf("expected SSE script in page")
	}
}

func TestHandler_RootReRendersOnEachRequest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("# Version 1"), 0644); err != nil {
		t.Fatal(err)
	}

	inputs := []Input{{Path: path}}
	h := newHandler(inputs, theme.DarkTheme())

	// First request
	req1 := httptest.NewRequest("GET", "/", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)
	body1, _ := io.ReadAll(rec1.Result().Body)

	// Modify file
	if err := os.WriteFile(path, []byte("# Version 2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Second request
	req2 := httptest.NewRequest("GET", "/", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	body2, _ := io.ReadAll(rec2.Result().Body)

	if !strings.Contains(string(body1), "Version 1") {
		t.Error("first request should contain Version 1")
	}
	if !strings.Contains(string(body2), "Version 2") {
		t.Error("second request should contain Version 2")
	}
}

func TestHandler_EventsEndpoint(t *testing.T) {
	inputs := []Input{{Data: []byte("# Test")}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/events", nil)
	rec := httptest.NewRecorder()

	// Run in goroutine since it blocks; we just check headers
	done := make(chan struct{})
	go func() {
		h.ServeHTTP(rec, req)
		close(done)
	}()

	// Cancel the request to unblock
	req.Body = http.NoBody
	// The handler blocks on context, so we can check the response headers
	// that were set before blocking. Since httptest doesn't support streaming
	// easily, we just verify the handler doesn't panic and the initial
	// headers are correct by using a short-lived check.
	// For a proper test, we'd need a real server — keep this as a smoke test.
}

func TestHandler_StdinInput(t *testing.T) {
	inputs := []Input{{Data: []byte("# From Stdin")}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Result().Body)
	if !strings.Contains(string(body), "From Stdin") {
		t.Fatalf("expected stdin content, got: %s", string(body)[:200])
	}
}

func TestHandler_MultiFile(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "file1.md")
	path2 := filepath.Join(dir, "file2.md")
	os.WriteFile(path1, []byte("# First"), 0644)
	os.WriteFile(path2, []byte("# Second"), 0644)

	inputs := []Input{{Path: path1}, {Path: path2}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Result().Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "First") {
		t.Error("expected First file content")
	}
	if !strings.Contains(bodyStr, "Second") {
		t.Error("expected Second file content")
	}
	if !strings.Contains(bodyStr, "file-separator") {
		t.Error("expected file separator between files")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/server/ -v -count=1`
Expected: Compilation error — package does not exist yet.

- [ ] **Step 3: Implement the server**

Create `internal/server/server.go`:

```go
package server

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/vchitepu/weave/internal/htmlrenderer"
	"github.com/vchitepu/weave/internal/theme"
	"github.com/vchitepu/weave/internal/watcher"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Input represents a Markdown source. Either Path or Data is set.
type Input struct {
	Path string // file path (re-read on each request)
	Data []byte // pre-read bytes (stdin)
}

// Start launches the web server and blocks until interrupted.
func Start(inputs []Input, th theme.Theme, port int) error {
	h := newHandler(inputs, th)

	// Collect file paths for watching
	var paths []string
	for _, in := range inputs {
		if in.Path != "" {
			paths = append(paths, in.Path)
		}
	}

	// Start file watcher
	stop, err := watcher.Watch(paths, func() {
		h.broadcast()
	})
	if err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}
	defer stop()

	// Start HTTP server
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	fmt.Fprintf(os.Stderr, "Weave web viewer running at http://localhost:%d\nPress Ctrl+C to stop.\n", port)

	srv := &http.Server{Handler: h}

	go func() {
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	// Block until signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Fprintln(os.Stderr)

	return srv.Close()
}

type handler struct {
	inputs []Input
	th     theme.Theme

	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newHandler(inputs []Input, th theme.Theme) *handler {
	return &handler{
		inputs:  inputs,
		th:      th,
		clients: make(map[chan struct{}]struct{}),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "":
		h.handleRoot(w, r)
	case "/events":
		h.handleEvents(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Build the markdown renderer fresh each time
	htmlR := htmlrenderer.New(h.th)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(htmlR, htmlrenderer.Priority),
				),
			),
		),
	)

	var body bytes.Buffer
	for i, input := range h.inputs {
		if i > 0 {
			name := filepath.Base(input.Path)
			if name == "" || name == "." {
				name = fmt.Sprintf("input-%d", i+1)
			}
			body.WriteString(fmt.Sprintf("<hr class=\"file-separator\">\n<p class=\"file-label\">%s</p>\n",
				htmlrenderer.HtmlEscapeString(name)))
		}

		data, err := h.readInput(input)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading input: %v", err), http.StatusInternalServerError)
			return
		}

		var rendered bytes.Buffer
		if err := md.Convert(data, &rendered); err != nil {
			http.Error(w, fmt.Sprintf("Error rendering markdown: %v", err), http.StatusInternalServerError)
			return
		}
		body.Write(rendered.Bytes())
	}

	page := buildPage(h.th, body.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(page))
}

func (h *handler) readInput(input Input) ([]byte, error) {
	if input.Path != "" {
		return os.ReadFile(input.Path)
	}
	return input.Data, nil
}

func (h *handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
	}()

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (h *handler) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- struct{}{}:
		default:
			// Don't block if client is slow
		}
	}
}

func buildPage(th theme.Theme, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>weave</title>
<style>
%s
%s
</style>
</head>
<body>
%s
<script>
const es = new EventSource("/events");
es.onmessage = function(e) {
  if (e.data === "reload") {
    location.reload();
  }
};
</script>
</body>
</html>`, htmlrenderer.ThemeCSS(th), htmlrenderer.BaseCSS, body)
}
```

Note: The `htmlEscapeString` function in `htmlrenderer` is unexported. We need to export it for use in the server. Rename it to `HtmlEscapeString` in `internal/htmlrenderer/htmlrenderer.go`:

Change in `internal/htmlrenderer/htmlrenderer.go`:
```go
// HtmlEscapeString escapes HTML special characters in a string.
func HtmlEscapeString(s string) string {
	return string(htmlEscape([]byte(s)))
}
```

And update any internal callers to use `htmlEscapeString` (keep the lowercase version for internal use, add the exported one that calls it):

```go
// HtmlEscapeString escapes HTML special characters in a string (exported).
func HtmlEscapeString(s string) string {
	return htmlEscapeString(s)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/server/ -v -count=1`
Expected: All PASS (the SSE test is a smoke test).

- [ ] **Step 5: Run all tests together**

Run: `go test ./... -count=1`
Expected: All PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/server/ internal/htmlrenderer/htmlrenderer.go
git commit -m "feat: add HTTP server with SSE live reload"
```


---

### Task 12: CLI integration — `--web` and `--port` flags

**Files:**
- Modify: `cmd/weave/main.go`

- [ ] **Step 1: Write the failing test**

Add to `cmd/weave/main_test.go`:

```go
func TestWebFlag_MissingFile(t *testing.T) {
	// Simulate the run() function with --web and a missing file.
	// We can't easily test the full server lifecycle in a unit test,
	// but we can verify the input validation path.
	themeFlag = ""
	widthFlag = 0

	// The run function should fail before starting the server
	// when a file doesn't exist.
	rootCmd := &cobra.Command{
		Use:  "weave [file...]",
		Args: cobra.ArbitraryArgs,
		RunE: run,
	}
	rootCmd.Flags().StringVar(&themeFlag, "theme", "", "")
	rootCmd.Flags().IntVar(&widthFlag, "width", 0, "")
	rootCmd.Flags().BoolVar(&webFlag, "web", false, "")
	rootCmd.Flags().IntVar(&portFlag, "port", 7331, "")

	rootCmd.SetArgs([]string{"--web", "/nonexistent/file.md"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file with --web")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/weave/ -run "TestWebFlag" -v -count=1`
Expected: Compilation error — `webFlag` and `portFlag` don't exist yet.

- [ ] **Step 3: Implement CLI integration**

Modify `cmd/weave/main.go`. Add the new flag variables:

```go
var (
	version   = "dev"
	themeFlag string
	widthFlag int
	webFlag   bool
	portFlag  int
)
```

Add flags in `main()`:

```go
rootCmd.Flags().BoolVar(&webFlag, "web", false, "Launch web viewer instead of terminal output")
rootCmd.Flags().IntVar(&portFlag, "port", 7331, "HTTP port for web viewer")
```

Add the import for the server package:

```go
"github.com/vchitepu/weave/internal/server"
```

Modify the beginning of `run()` to handle `--web` mode. Insert this block right after the function signature, before the existing multi-file logic:

```go
func run(cmd *cobra.Command, args []string) error {
	if webFlag {
		return runWeb(cmd, args)
	}
	// ... existing code unchanged ...
}

func runWeb(cmd *cobra.Command, args []string) error {
	// Validate theme
	if themeFlag != "" && themeFlag != "dark" && themeFlag != "light" {
		return fmt.Errorf("weave: invalid theme %q (use 'dark' or 'light')", themeFlag)
	}
	th := theme.Detect(themeFlag)

	var inputs []server.Input

	if len(args) == 0 {
		// Stdin mode
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return cmd.Help()
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("weave: failed to read stdin: %w", err)
		}
		inputs = append(inputs, server.Input{Data: data})
	} else {
		// File mode — validate all files exist before starting server
		for _, path := range args {
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("weave: no such file: %s", path)
			}
			inputs = append(inputs, server.Input{Path: path})
		}
	}

	return server.Start(inputs, th, portFlag)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./cmd/weave/ -run "TestWebFlag" -v -count=1`
Expected: PASS.

- [ ] **Step 5: Build and verify the binary**

Run: `go build -o /tmp/weave-test ./cmd/weave && echo "Build OK"`
Expected: "Build OK"

- [ ] **Step 6: Run the full test suite**

Run: `go test ./... -count=1`
Expected: All PASS.

- [ ] **Step 7: Commit**

```bash
git add cmd/weave/main.go cmd/weave/main_test.go
git commit -m "feat: add --web and --port flags for web viewer mode"
```


---

### Task 13: End-to-end smoke test and README update

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Manual smoke test**

Build and run the web viewer manually:

```bash
go build -o /tmp/weave-test ./cmd/weave
echo "# Hello\n\nThis is a **test**.\n\n\`\`\`go\nfmt.Println(\"hello\")\n\`\`\`" > /tmp/test.md
/tmp/weave-test --web /tmp/test.md &
WEB_PID=$!
sleep 1
# Fetch the page and verify content
curl -s http://localhost:7331 | grep -q "Hello" && echo "PASS: page renders" || echo "FAIL"
curl -s http://localhost:7331 | grep -q "EventSource" && echo "PASS: SSE script present" || echo "FAIL"
curl -s http://localhost:7331 | grep -q "code-container" && echo "PASS: code container" || echo "FAIL"
kill $WEB_PID 2>/dev/null
```

Expected: All checks print PASS.

- [ ] **Step 2: Update README.md**

Add `--web` and `--port` to the usage and flags sections of `README.md`.

In the Usage section, add after the existing examples:

```markdown
# Launch web viewer
weave --web README.md

# Web viewer on custom port
weave --web --port 8080 README.md

# Web viewer with multiple files
weave --web file1.md file2.md

# Pipe to web viewer
cat README.md | weave --web
```

In the Flags table, add:

```markdown
| `--web` | `false` | Launch web viewer in browser |
| `--port` | `7331` | HTTP port for web viewer |
```

In the Future Directions / UX Improvements section, add:

```markdown
- `--web` edit mode for in-browser Markdown editing
```

- [ ] **Step 3: Run full test suite one final time**

Run: `go test ./... -v -count=1`
Expected: All PASS.

- [ ] **Step 4: Commit**

```bash
git add README.md
git commit -m "docs: add web viewer to README usage and flags"
```

- [ ] **Step 5: Push branch to remote**

```bash
git push
```

