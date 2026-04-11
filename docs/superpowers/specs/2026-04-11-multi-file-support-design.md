# Multi-File Support Design

**Date:** 2026-04-11
**Branch:** feature/multi-file-support

## Overview

Allow `weave` to accept multiple file arguments and render them sequentially in a single output stream, separated by a visual section break that identifies each file.

```sh
weave file1.md file2.md file3.md
```

## Goals

- Render 2+ Markdown files in one invocation with clear visual separators
- Reuse existing renderer, theme, width, and paging infrastructure unchanged
- Fail fast if any file cannot be read

## Non-Goals

- Stdin + file mixing (not supported; stdin is ignored when 2+ file args are given)
- Glob expansion (the shell handles this; `weave *.md` already works via shell expansion)
- Per-file theme or width overrides

## CLI Changes

**File:** `cmd/weave/main.go`

- Change `cobra.MaximumNArgs(1)` to `cobra.ArbitraryArgs`
- When `len(args) == 0`: existing stdin path (unchanged)
- When `len(args) == 1`: existing single-file path (unchanged)
- When `len(args) >= 2`: multi-file path (new)

`--theme` and `--width` flags apply globally to all files.

## `renderFile()` Helper

A new unexported function added to `cmd/weave/main.go`:

```go
func renderFile(path string, md goldmark.Markdown) (string, error)
```

- Reads the file with `os.ReadFile(path)`
- Renders via `md.Convert(input, &buf)`
- Returns the rendered string on success
- Returns an error immediately on read failure (fail fast)
- The `goldmark.Markdown` instance is constructed once in `run()` with theme and width already baked in, then passed to each `renderFile()` call

## Multi-File Loop

In `run()`, when `len(args) >= 2`:

1. Construct the `goldmark.Markdown` instance once (same as today)
2. Initialise a `strings.Builder` to accumulate all output
3. For each path in `args`:
   - If not the first file, append a separator (see below)
   - Call `renderFile(path, md)` — return the error immediately on failure
   - Append the rendered output to the builder
4. Make the single paging decision on the total accumulated output (same logic as today)

## Section Separator

Between consecutive files, the separator consists of two lines:

1. A full-width thematic break rule (`────────...`) using the existing thematic break style from the renderer
2. The filename on its own line, styled as a dim/muted label (not a heading level)

The separator is rendered as a plain string (not passed through goldmark) by a small `fileSeparator(filename string, width int, th theme.Theme) string` helper in `main.go`. It reuses the existing lipgloss styling primitives already used elsewhere.

## Error Handling

- Any file that cannot be read causes `run()` to return that error immediately
- Cobra prints the error to stderr and exits non-zero
- Because output is buffered until all rendering is complete, nothing is written to stdout when an error occurs — the user sees only the error message, not partial output

## Testing

**`cmd/weave/main_test.go`** additions:

| Test | What it checks |
|------|----------------|
| `TestRenderFile_ValidFile` | `renderFile()` with a valid temp `.md` file returns non-empty output and no error |
| `TestRenderFile_MissingFile` | `renderFile()` with a non-existent path returns a non-nil error |
| `TestFileSeparator` | `fileSeparator()` output contains the filename and at least one `─` character |
| `TestMultiFileOutput` | Two temp `.md` files rendered together: output contains both filenames and a separator `─` between them |

Existing tests remain unaffected.

## File Impact Summary

| File | Change |
|------|--------|
| `cmd/weave/main.go` | `ArbitraryArgs`, multi-file loop, `renderFile()`, `fileSeparator()` |
| `cmd/weave/main_test.go` | 4 new tests |
| All other files | No changes |
