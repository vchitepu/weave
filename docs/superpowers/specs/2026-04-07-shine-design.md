# Shine — Terminal Markdown Viewer: Design Spec

**Date:** 2026-04-07  
**Status:** Approved  
**Language:** Go  

## Overview

`shine` is a terminal Markdown viewer written in Go. It aims to capture the visual richness of web-based Markdown renderers — explicit styled containers for code blocks, visual wrappers for ASCII diagrams and tree output, syntax highlighting, and clean table rendering — while remaining a fast, composable Unix CLI tool.

It is inspired by `glow` but built on a custom AST renderer rather than `glamour`, giving full control over every visual element.

---

## Architecture

The tool is structured as a rendering pipeline:

```
Input (stdin or file)
       ↓
   Markdown source ([]byte)
       ↓
   goldmark parser → AST
       ↓
   shine AST walker/renderer
       ↓
   lipgloss-styled strings
       ↓
   width-aware output writer
       ↓
   stdout (or $PAGER if overflow)
```

### Package Structure

```
cmd/shine/          CLI entry point (cobra)
internal/renderer/  goldmark Renderer implementation
internal/theme/     Theme definitions and auto-detection
```

---

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/yuin/goldmark` | Markdown parser (AST) |
| `github.com/alecthomas/chroma` | Syntax highlighting inside code blocks |
| `github.com/charmbracelet/lipgloss` | Styled terminal containers, borders, colors |
| `github.com/spf13/cobra` | CLI flag/argument handling |
| `golang.org/x/term` | Terminal width/height detection, TTY check |

---

## CLI Interface

Built with `cobra`. No subcommands.

### Usage

```
shine [flags] [file]

shine README.md
cat README.md | shine
some-cmd --help | shine
shine                     # reads from stdin, waits for EOF
```

### Flags

```
--theme dark|light    Override auto-detected theme
--width N             Override terminal width (default: auto-detect)
--version             Print version and exit
--help                Print usage and exit
```

### Paging

- Render full output to a buffer first.
- Detect terminal height via `golang.org/x/term`.
- If rendered line count exceeds terminal height, exec `$PAGER` (fallback: `less -R`).
- If stdout is not a TTY, skip paging and write raw output (enables `shine README.md | grep foo`).

### Error Handling

- File not found → `shine: no such file: foo.md` to stderr, exit 1
- Unreadable input (no TTY, no pipe) → print usage, exit 1
- Unknown flag → print usage, exit 1

---

## Rendering Elements

### Headings

- **H1:** Bold + themed color, full-width `─` rule beneath
- **H2:** Bold + themed color, no rule
- **H3–H6:** Progressively dimmer via color only (no bold below H3)

### Code Blocks — Standard

Rounded `lipgloss` border, full terminal width. Language badge in top-left of border. Content syntax-highlighted via `chroma`.

```
╭─ go ──────────────────────────────────────────╮
│ func main() {                                  │
│     fmt.Println("hello")                       │
│ }                                              │
╰────────────────────────────────────────────────╯
```

### Inline Code

Subtle background highlight, no border.

### Blockquotes

Left bar `│` in accent color, indented content.

### Lists

- Unordered: `•` bullet, nested lists indent with two spaces
- Ordered: `1.` numbers, same nesting

### Tables

ASCII box-drawing borders, column-aligned, header row bold + colored.

### Horizontal Rules

Full-width `─` line in dim color.

### Images

Rendered as `[image: alt text]` in italic dim style. No inline image rendering in v1.

### Bold / Italic / Strikethrough

Mapped to ANSI bold, italic, strikethrough respectively.

### Links

Link text rendered normally, URL appended in dim: `text (url)`.

---

## Visual Container Detection

Detection is based solely on the info string (language hint) of fenced code blocks. No content heuristics in v1.

| Info string | Container type | Border color |
|---|---|---|
| `tree` | Directory tree | Green |
| `ascii`, `diagram`, `art` | ASCII diagram | Blue |
| `mermaid` | Diagram (not rendered) | Blue |
| `bash`, `sh`, `shell`, `console`, `terminal` | Shell session | Yellow |
| anything else | Standard code | Default (theme accent) |

### Container Examples

**Tree:**
```
╭─ tree ─────────────────────────────────────────╮
│ src/                                            │
│ ├── main.go                                     │
│ └── renderer/                                   │
╰─────────────────────────────────────────────────╯
```

**Diagram:**
```
╭─ diagram ──────────────────────────────────────╮
│   ┌──────┐     ┌──────┐                        │
│   │  A   │────▶│  B   │                        │
│   └──────┘     └──────┘                        │
╰────────────────────────────────────────────────╯
```

**Shell:**
```
╭─ $ ────────────────────────────────────────────╮
│ $ go build ./...                               │
│ $ ./shine README.md                            │
╰────────────────────────────────────────────────╯
```

**Mermaid (not rendered):**
```
╭─ diagram ──────────────────────────────────────╮
│ [diagram not rendered]                         │
│                                                │
│ graph TD                                       │
│     A --> B                                    │
╰────────────────────────────────────────────────╯
```

---

## Theme System

### Auto-Detection Order

1. `$SHINE_THEME` env var (`dark` or `light`) — explicit override
2. `$COLORFGBG` — format `fg;bg`; background luma < 128 = dark
3. `$TERM_PROGRAM` heuristics — Apple Terminal defaults light, most others dark
4. Fallback: dark

### Theme Struct

```go
type Theme struct {
    // Text
    Normal        lipgloss.Style
    Bold          lipgloss.Style
    Italic        lipgloss.Style
    Strikethrough lipgloss.Style
    Dim           lipgloss.Style

    // Headings
    H1, H2, H3  lipgloss.Style
    HeadingRule  lipgloss.Style

    // Code
    CodeBorder   lipgloss.Color
    CodeHeader   lipgloss.Style
    InlineCode   lipgloss.Style

    // Container variants
    TreeBorder    lipgloss.Color  // green
    DiagramBorder lipgloss.Color  // blue
    ShellBorder   lipgloss.Color  // yellow

    // Blockquote
    BlockquoteBar lipgloss.Style

    // Table
    TableHeader  lipgloss.Style
    TableBorder  lipgloss.Color

    // Links / Images
    LinkURL      lipgloss.Style
    ImageAlt     lipgloss.Style
}
```

Two concrete instances: `DarkTheme` and `LightTheme` — plain Go vars. No config files or JSON themes in v1.

---

## Out of Scope (v1)

- Inline image rendering (iTerm2/Kitty protocols)
- Mermaid diagram rendering
- User-configurable theme files
- Interactive pager (built-in scroll)
- Watch mode (`--watch` for live reload)
- Multiple file arguments


---

## Implementation Notes

- `Theme.H1`, `Theme.H2`, `Theme.H3` are the three distinct heading styles. H4–H6 reuse `H3` style with no further distinction in v1.
- `cobra` provides `--help` and `--version` as builtins; `--theme` and `--width` are registered as persistent flags on the root command.
- The renderer passes `terminalWidth` (int) at construction time; all containers and rules are sized relative to it.
