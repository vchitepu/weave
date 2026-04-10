# Design: Task List Checkboxes

**Date:** 2026-04-10
**Branch:** task-list-checkboxes

## Overview

Add support for GitHub-Flavored Markdown task list items (`- [x]` / `- [ ]`) to weave's terminal renderer. Checked and unchecked items render with distinct themed symbols, giving users a clear visual distinction between done and pending tasks.

## Approach

Use goldmark's built-in `extension.TaskList`. This extension parses `- [x]` and `- [ ]` syntax and emits `ast.TaskCheckBox` nodes into the parse tree. The renderer registers a handler for this node type, rendering the checkbox symbol before the list item text. This is idiomatic — the same pattern used for `extension.Table` and `extension.Strikethrough`.

## Components

### `cmd/weave/main.go`

Add `extension.TaskList` to the goldmark `WithExtensions(...)` call. One line change alongside the existing `extension.Table` and `extension.Strikethrough`.

### `internal/theme/theme.go`

Add two new `lipgloss.Style` fields to the `Theme` struct:

- `TaskChecked` — green foreground (or success-toned), used for `✓`
- `TaskUnchecked` — dimmed foreground, used for `○`

Wire both fields into `DarkTheme()` and `LightTheme()`. Dark theme uses a muted green for checked and a dim gray for unchecked. Light theme mirrors this with appropriate contrast-adjusted colors.

### `internal/renderer/list.go`

Register a `renderTaskCheckBox` handler in `RegisterFuncs`. On entering a `TaskCheckBox` node:

- If `node.IsChecked`: write `r.theme.TaskChecked.Render("✓ ")`
- Otherwise: write `r.theme.TaskUnchecked.Render("○ ")`

Write directly to the output buffer (not `paraBuf`).

**Bullet suppression**: The `renderListItem` handler writes the bullet prefix (`•` or `1.`) before walking children. For task list items the bullet must be suppressed so the output is `  ✓ Done` rather than `  • ✓ Done`. Detection: check whether `node.FirstChild()` is an `ast.TaskCheckBox` (goldmark's extension always inserts it as the first child of a task `ListItem`). If so, write only the indent (no bullet), and let `renderTaskCheckBox` write the symbol. The `listPrefixWidths` entry should account for the indent width so word-wrap continuation lines align correctly.

### `internal/renderer/list_test.go`

Add two test cases using the shared `renderMarkdown` helper:

1. `- [x] Done item` — assert output contains `✓` and `Done item`
2. `- [ ] Pending item` — assert output contains `○` and `Pending item`
3. Mixed list — assert both symbols appear and items are distinct

## Data Flow

```
Markdown input: "- [x] Done\n- [ ] Todo"
  → goldmark parser + extension.TaskList
  → AST: List > ListItem > TaskCheckBox(IsChecked=true) + Text("Done")
          List > ListItem > TaskCheckBox(IsChecked=false) + Text("Todo")
  → renderListItem: writes indent + bullet prefix (suppressed for task items?)
  → renderTaskCheckBox: writes themed ✓ or ○
  → renderText: writes item text
  → output: "  ✓ Done\n  ○ Todo"
```

Note: goldmark's TaskList extension replaces the `[ ]`/`[x]` text from the parsed content, so the raw checkbox syntax won't appear in the output. The `TaskCheckBox` node is the first child of a task `ListItem` and renders before the item's text children. The bullet suppression logic in `renderListItem` ensures no double prefix.

## Out of Scope

- Toggling checkboxes interactively (weave is a viewer, not an editor)
- Nested task lists (handled naturally by existing list nesting logic)
- `- [X]` capital X — goldmark's extension handles this already
