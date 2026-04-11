# weave

A terminal Markdown viewer with rich visual containers. Renders Markdown with styled code blocks, tables, blockquotes, lists, and headings — directly in your terminal.

## Features

- Syntax-highlighted fenced code blocks with language badges
- Automatic detection of `tree`, `diagram`, and `shell` code block variants
- Tables with box-drawing borders and column alignment
- Blockquotes with bar markers, word-wrapped to terminal width
- Ordered and unordered lists with nested indentation
- Task list checkboxes with checked/unchecked symbols
- Dark and light themes with auto-detection
- Auto-paging via `$PAGER` when output exceeds terminal height
- Pipe-friendly — skips paging when stdout is not a TTY

## Installation

Requires Go 1.21 or later.

```sh
go install github.com/vchitepu/weave/cmd/weave@latest
```

Or build from source:

```sh
git clone https://github.com/vchitepu/weave.git
cd weave
go build -o ~/go/bin/weave ./cmd/weave
```

## Usage

```sh
# Render a file
weave README.md

# Pipe from stdin
cat README.md | weave

# Pipe from a command
some-cmd --help | weave

# Force a theme
weave --theme dark README.md
weave --theme light README.md

# Override terminal width
weave --width 100 README.md
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--theme` | auto-detected | Theme to use: `dark` or `light` |
| `--width` | auto-detected | Override terminal width (cols) |
| `--version` | | Print version and exit |
| `--help` | | Print usage and exit |

## Theme Auto-Detection

weave detects your terminal theme in this order:

1. `--theme` flag
2. `WEAVE_THEME` environment variable (`dark` or `light`)
3. `COLORFGBG` environment variable
4. Terminal program heuristics (`Apple_Terminal` → light)
5. Falls back to `dark`

## License

MIT
